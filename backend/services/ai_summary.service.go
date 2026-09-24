package services

import (
	"backend-queue/config"
	"backend-queue/dto"
	"backend-queue/repositories"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type AISummaryService struct {
	queueRepo   *repositories.QueueRepository
	summaryRepo *repositories.SummaryRepository
}

func NewAISummaryService() *AISummaryService {
	return &AISummaryService{
		queueRepo:   repositories.NewQueueRepository(),
		summaryRepo: repositories.NewSummaryRepository(),
	}
}

// GetDailySummary returns an AI-generated summary for the given date.
// Results are cached in Redis for 5 minutes to avoid spamming Gemini API.
func (s *AISummaryService) GetDailySummary(date string) (*dto.AISummaryResponse, error) {
	ctx := context.Background()
	cacheKey := "dashboard:ai_summary:" + date

	// ── 1. Try Redis cache first ─────────────────────────────────────────────
	cached, err := config.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp dto.AISummaryResponse
		if jsonErr := json.Unmarshal([]byte(cached), &resp); jsonErr == nil {
			resp.Cached = true
			return &resp, nil
		}
	}

	// ── 2. Cache miss — collect data from DB ─────────────────────────────────
	stats, err := s.queueRepo.GetStatsByDate(date)
	if err != nil {
		return nil, fmt.Errorf("GetDailySummary stats: %w", err)
	}

	breakdown, err := s.summaryRepo.GetServiceBreakdown(date)
	if err != nil {
		return nil, fmt.Errorf("GetDailySummary breakdown: %w", err)
	}

	peakHour, peakCount, err := s.summaryRepo.GetPeakHour(date)
	if err != nil {
		return nil, fmt.Errorf("GetDailySummary peak: %w", err)
	}

	avgMinutes, err := s.queueRepo.GetAvgServiceMinutes(date)
	if err != nil {
		return nil, fmt.Errorf("GetDailySummary avg: %w", err)
	}

	// ── 3. Build service breakdown string ────────────────────────────────────
	var serviceLines []string
	for _, b := range breakdown {
		serviceLines = append(serviceLines, fmt.Sprintf("  - %s: %d antrian", b.Name, b.Count))
	}
	serviceText := "  (tidak ada data)"
	if len(serviceLines) > 0 {
		serviceText = strings.Join(serviceLines, "\n")
	}

	avgText := "tidak ada data selesai"
	if avgMinutes > 0 {
		avgText = fmt.Sprintf("%.1f menit", avgMinutes)
	}

	peakText := "tidak ada data"
	if peakCount > 0 {
		peakText = fmt.Sprintf("Pukul %02d.00 (%d antrian masuk)", peakHour, peakCount)
	}

	// ── 4. Build Gemini prompt ────────────────────────────────────────────────
	prompt := fmt.Sprintf(`Kamu adalah asisten analitik showroom otomotif.
Berikan ringkasan dan insight dari data antrian berikut untuk tanggal %s:

Statistik Antrian:
- Total antrian: %d
- Menunggu: %d
- Sedang diproses: %d
- Selesai: %d
- Dibatalkan: %d

Rata-rata waktu servis aktual: %s
Peak hour (jam tersibuk): %s

Breakdown per jenis layanan:
%s

Berikan respons dalam format berikut (gunakan tanda newline antar bagian):
1. RINGKASAN SINGKAT: (2-3 kalimat ringkasan operasional hari ini)
2. PEAK HOUR: (analisis jam tersibuk dan dampaknya)
3. REKOMENDASI OPERASIONAL: (2-3 rekomendasi konkret dan actionable untuk manajemen showroom)

Gunakan bahasa Indonesia yang profesional dan ringkas. Jangan gunakan markdown bold/italic.`,
		date,
		stats["total"],
		stats["waiting"],
		stats["processing"],
		stats["done"],
		stats["cancelled"],
		avgText,
		peakText,
		serviceText,
	)

	// ── 5. Call Gemini REST API ───────────────────────────────────────────────
	summaryText, err := callGeminiAPI(prompt)
	if err != nil {
		return nil, fmt.Errorf("Gemini API error: %w", err)
	}

	// ── 6. Build response ─────────────────────────────────────────────────────
	resp := &dto.AISummaryResponse{
		Summary:     summaryText,
		GeneratedAt: time.Now().Format(time.RFC3339),
		Cached:      false,
	}

	// ── 7. Save to Redis with 5-minute TTL ───────────────────────────────────
	if data, jsonErr := json.Marshal(resp); jsonErr == nil {
		config.RDB.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return resp, nil
}

// ── Gemini REST API helper ────────────────────────────────────────────────────

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func callGeminiAPI(prompt string) (string, error) {
	apiKey := config.GetEnv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY tidak ditemukan di environment")
	}

	model := config.GetEnv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-3.6-flash"
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		model,
		apiKey,
	)

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpResp, err := http.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("HTTP request gagal: %w", err)
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "", fmt.Errorf("baca response: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API status %d: %s", httpResp.StatusCode, string(respBytes))
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(respBytes, &geminiResp); err != nil {
		return "", fmt.Errorf("parse Gemini response: %w", err)
	}

	if geminiResp.Error != nil {
		return "", fmt.Errorf("Gemini error: %s", geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("Gemini tidak mengembalikan konten")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
