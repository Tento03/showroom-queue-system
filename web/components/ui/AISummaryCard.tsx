'use client'

import { useState } from 'react'
import { getDailySummary } from '@/lib/api'
import { AISummaryResponse } from '@/lib/types'
import { format } from 'date-fns'
import { id as localeId } from 'date-fns/locale'

interface AISummaryCardProps {
  date: string
}

function SpinnerIcon({ size = 16 }: { size?: number }) {
  return (
    <svg
      className="spinner"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={2.5}
    >
      <circle cx="12" cy="12" r="10" strokeOpacity={0.25} />
      <path d="M22 12a10 10 0 0 0-10-10" />
    </svg>
  )
}

export default function AISummaryCard({ date }: AISummaryCardProps) {
  const [summary, setSummary] = useState<AISummaryResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleGenerate = async () => {
    setLoading(true)
    setError(null)
    try {
      const result = await getDailySummary(date)
      setSummary(result)
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Terjadi kesalahan saat generate summary'
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  // Reset when date changes (handled by parent passing new date prop)
  // The summary remains visible until explicitly regenerated

  return (
    <div
      style={{
        background: 'linear-gradient(135deg, #1e1b4b 0%, #312e81 40%, #4338ca 100%)',
        borderRadius: '1rem',
        padding: '1.5rem',
        boxShadow: '0 4px 24px rgba(67,56,202,0.25)',
        position: 'relative',
        overflow: 'hidden',
      }}
    >
      {/* Decorative glow blob */}
      <div
        style={{
          position: 'absolute',
          top: '-40px',
          right: '-40px',
          width: '160px',
          height: '160px',
          borderRadius: '50%',
          background: 'rgba(139,92,246,0.25)',
          filter: 'blur(40px)',
          pointerEvents: 'none',
        }}
      />

      {/* Header */}
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '1rem' }}>
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <span style={{ fontSize: '1.25rem' }}>🤖</span>
            <h2 style={{ color: '#e0e7ff', fontWeight: 700, fontSize: '1rem', margin: 0 }}>
              AI Insight Harian
            </h2>
          </div>
          <p style={{ color: '#a5b4fc', fontSize: '0.75rem', marginTop: '0.25rem' }}>
            Powered by Gemini AI · {date}
          </p>
        </div>

        <button
          onClick={handleGenerate}
          disabled={loading}
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: '0.5rem',
            padding: '0.5rem 1.25rem',
            borderRadius: '0.75rem',
            border: 'none',
            background: loading ? 'rgba(255,255,255,0.1)' : 'linear-gradient(135deg, #7c3aed, #6366f1)',
            color: '#fff',
            fontWeight: 600,
            fontSize: '0.8125rem',
            cursor: loading ? 'not-allowed' : 'pointer',
            opacity: loading ? 0.7 : 1,
            transition: 'all 0.2s',
            boxShadow: loading ? 'none' : '0 2px 8px rgba(124,58,237,0.4)',
            whiteSpace: 'nowrap',
          }}
          onMouseEnter={e => {
            if (!loading) (e.currentTarget as HTMLButtonElement).style.transform = 'translateY(-1px)'
          }}
          onMouseLeave={e => {
            (e.currentTarget as HTMLButtonElement).style.transform = 'translateY(0)'
          }}
        >
          {loading ? <SpinnerIcon size={14} /> : <span>✨</span>}
          {loading ? 'Menganalisis...' : 'Generate Summary'}
        </button>
      </div>

      {/* Content area */}
      {!summary && !loading && !error && (
        <div
          style={{
            border: '1px dashed rgba(165,180,252,0.3)',
            borderRadius: '0.75rem',
            padding: '2rem',
            textAlign: 'center',
          }}
        >
          <p style={{ color: '#818cf8', fontSize: '0.875rem', margin: 0 }}>
            Klik <strong style={{ color: '#a5b4fc' }}>&quot;Generate Summary&quot;</strong> untuk mendapatkan
            insight analitik AI dari data antrian hari ini.
          </p>
        </div>
      )}

      {loading && (
        <div
          style={{
            border: '1px dashed rgba(165,180,252,0.3)',
            borderRadius: '0.75rem',
            padding: '2rem',
            textAlign: 'center',
          }}
        >
          <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '0.75rem' }}>
            <SpinnerIcon size={24} />
            <p style={{ color: '#a5b4fc', fontSize: '0.875rem', margin: 0 }}>
              Menganalisis data antrian dan meminta insight dari Gemini AI...
            </p>
          </div>
        </div>
      )}

      {error && (
        <div
          style={{
            background: 'rgba(239,68,68,0.15)',
            border: '1px solid rgba(239,68,68,0.3)',
            borderRadius: '0.75rem',
            padding: '1rem',
            display: 'flex',
            alignItems: 'flex-start',
            gap: '0.5rem',
          }}
        >
          <span style={{ fontSize: '1rem', flexShrink: 0 }}>⚠️</span>
          <div>
            <p style={{ color: '#fca5a5', fontWeight: 600, fontSize: '0.875rem', margin: '0 0 0.25rem' }}>
              Gagal generate summary
            </p>
            <p style={{ color: '#fca5a5', fontSize: '0.8125rem', margin: 0, opacity: 0.8 }}>{error}</p>
          </div>
        </div>
      )}

      {summary && !loading && (
        <div
          style={{
            background: 'rgba(255,255,255,0.06)',
            backdropFilter: 'blur(8px)',
            borderRadius: '0.75rem',
            padding: '1.25rem',
            border: '1px solid rgba(165,180,252,0.2)',
            animation: 'fadeIn 0.4s ease',
          }}
        >
          {/* Summary text */}
          <p
            style={{
              color: '#e0e7ff',
              fontSize: '0.875rem',
              lineHeight: '1.75',
              whiteSpace: 'pre-wrap',
              margin: '0 0 1rem',
            }}
          >
            {summary.summary}
          </p>

          {/* Footer metadata */}
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '0.75rem',
              paddingTop: '0.75rem',
              borderTop: '1px solid rgba(165,180,252,0.15)',
            }}
          >
            <span style={{ color: '#6366f1', fontSize: '0.75rem' }}>
              🕐 Dibuat:{' '}
              {format(new Date(summary.generated_at), "dd MMM yyyy, HH:mm", { locale: localeId })}
            </span>

            {summary.cached && (
              <span
                style={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  gap: '0.25rem',
                  padding: '0.125rem 0.625rem',
                  borderRadius: '9999px',
                  background: 'rgba(99,102,241,0.2)',
                  border: '1px solid rgba(99,102,241,0.3)',
                  color: '#a5b4fc',
                  fontSize: '0.6875rem',
                  fontWeight: 600,
                }}
              >
                📦 Cached
              </span>
            )}

            <button
              onClick={handleGenerate}
              style={{
                marginLeft: 'auto',
                background: 'transparent',
                border: '1px solid rgba(165,180,252,0.3)',
                borderRadius: '0.5rem',
                color: '#a5b4fc',
                fontSize: '0.75rem',
                padding: '0.25rem 0.75rem',
                cursor: 'pointer',
                transition: 'all 0.15s',
              }}
              onMouseEnter={e => {
                (e.currentTarget as HTMLButtonElement).style.background = 'rgba(165,180,252,0.1)'
              }}
              onMouseLeave={e => {
                (e.currentTarget as HTMLButtonElement).style.background = 'transparent'
              }}
            >
              🔄 Refresh
            </button>
          </div>
        </div>
      )}

      <style>{`
        @keyframes fadeIn {
          from { opacity: 0; transform: translateY(6px); }
          to   { opacity: 1; transform: translateY(0); }
        }
        .spinner {
          animation: spin 0.8s linear infinite;
          transform-origin: center;
        }
        @keyframes spin {
          from { transform: rotate(0deg); }
          to   { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  )
}
