# 🚗 Showroom Queue System

A real-time vehicle queue management system built for automotive showrooms. Staff can register vehicles via mobile OCR scanning, while admins monitor and manage queues through a live web dashboard.

**Live Demo:** [showroom-queue-system-beta.vercel.app](https://showroom-queue-system-beta.vercel.app/dashboard)

![Dashboard](docs/dashboard.png)

---

## ✨ Features

- 📷 **OCR Plate Recognition** — scan vehicle plates automatically via mobile camera
- ⚡ **Real-time Dashboard** — queue updates instantly via WebSocket, no refresh needed
- 🔢 **Auto Queue Numbering** — race-condition-safe ticket generation using DB transactions
- 📊 **Live Statistics** — total, waiting, processing, done, cancelled — cached with Redis
- 🔄 **Status Management** — controlled transitions (waiting → processing → done)
- 🛡️ **Rate Limiting** — upload endpoint protected via Redis (5 req/10s per IP)
- 📄 **Pagination** — efficient data loading with page & limit support
- 🔍 **Search & Filter** — search by plate number or owner name, filter by date

---

## 🏗️ Architecture

```
[Flutter App - Staff]          [Next.js Web - Admin]
        │                               │
        │ HTTP REST                     │ HTTP + WebSocket
        └──────────────┬────────────────┘
                       ▼
              [Go Gin Backend]
                   │      │      │
                MySQL   Redis  WebSocket Hub
              (data)  (cache)  (realtime)
```

---

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| Mobile | Flutter + Google ML Kit (OCR) |
| Backend | Go (Gin, GORM) |
| Web | Next.js + Tailwind CSS |
| Database | MySQL |
| Cache & Rate Limit | Redis |
| Realtime | WebSocket |
| Deploy Backend | Railway |
| Deploy Web | Vercel |

---

## 🔑 API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/queue` | Register new vehicle (race-condition safe) |
| `GET` | `/queues` | List queues with pagination & date filter |
| `GET` | `/queue/:id` | Get queue detail |
| `PATCH` | `/queue/:id/status` | Update queue status |
| `DELETE` | `/queue/:id` | Delete queue |
| `GET` | `/dashboard/stats` | Get stats (Redis cached) |
| `POST` | `/upload` | Upload vehicle image (rate limited) |
| `GET` | `/ws` | WebSocket connection |

---

## ⚙️ Key Technical Decisions

### Race Condition Fix
Queue number generation is wrapped in a DB transaction with row-level locking (`SELECT FOR UPDATE`), preventing duplicate ticket numbers under concurrent requests.

### Redis Caching
Dashboard stats are cached in Redis with a 30-second TTL for today's data and 24-hour TTL for historical data. Cache is invalidated immediately on queue creation or status update.

### WebSocket Broadcast
A central hub manages all connected dashboard clients. Every queue event (created/updated) is broadcast instantly to all connected admins without polling.

### Status Transition Rules
```
waiting → processing → done
waiting → cancelled
processing → cancelled
```
Done and cancelled are terminal states — no further transitions allowed.

---

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- MySQL 8+
- Redis 7+
- Node.js 18+
- Flutter 3+

### Backend
```bash
cd backend
cp .env.example .env
# Fill in your DB and Redis credentials
go mod download
go run main.go
```

### Web Dashboard
```bash
cd web
cp .env.local.example .env.local
# Fill in API URL
npm install
npm run dev
```

### Mobile
```bash
cd mobile/vehicle_queue_app
flutter pub get
flutter run
```

---

## 📁 Project Structure

```
showroom-queue-system/
├── backend/
│   ├── config/          # DB, Redis, Env
│   ├── controllers/     # HTTP handlers
│   ├── services/        # Business logic + WebSocket hub
│   ├── repositories/    # DB queries
│   ├── middleware/       # CORS, Rate limit
│   ├── models/          # GORM models
│   ├── dto/             # Request/Response types
│   └── utils/           # Typed errors
├── web/
│   ├── app/             # Next.js App Router
│   ├── components/      # Reusable UI components
│   └── lib/             # API client, Types
└── mobile/
    └── vehicle_queue_app/
        └── lib/
            ├── features/ # Feature-first structure
            └── core/     # Config, shared
```

---

## 🌐 Deployment

| Service | Platform | URL |
|---|---|---|
| Backend + DB + Cache | Railway | `showroom-queue-system-production.up.railway.app` |
| Web Dashboard | Vercel | `showroom-queue-system-beta.vercel.app` |

---

## 👨‍💻 Author

**Tento** — Full-stack Developer  
[GitHub](https://github.com/Tento03)