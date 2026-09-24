'use client'

import { useEffect, useState, useCallback, useRef } from 'react'
import { format, isBefore } from 'date-fns'
import { Queue, DashboardStats, QueueStatus, ETAResponse, QueueETA } from '@/lib/types'
import { getQueues, getDashboardStats, updateQueueStatus, deleteQueue, getEstimates } from '@/lib/api'
import StatsCard from '@/components/ui/StatsCard'
import AISummaryCard from '@/components/ui/AISummaryCard'

// ─── Status Badge ──────────────────────────────────────────────────────────────

const STATUS_CONFIG: Record<string, { label: string; className: string; dot?: string }> = {
  waiting: {
    label: 'Menunggu',
    className: 'bg-amber-50 text-amber-700 border border-amber-200',
    dot: 'bg-amber-400',
  },
  processing: {
    label: 'Diproses',
    className: 'bg-blue-50 text-blue-700 border border-blue-200 pulse-processing',
    dot: 'bg-blue-500',
  },
  done: {
    label: 'Selesai',
    className: 'bg-emerald-50 text-emerald-700 border border-emerald-200',
    dot: 'bg-emerald-500',
  },
  cancelled: {
    label: 'Batal',
    className: 'bg-red-50 text-red-600 border border-red-200',
    dot: 'bg-red-400',
  },
}

function StatusBadge({ status }: { status: string }) {
  const cfg = STATUS_CONFIG[status] ?? { label: status, className: 'bg-gray-100 text-gray-600', dot: 'bg-gray-400' }
  return (
    <span className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-semibold transition-all duration-300 ${cfg.className}`}>
      <span className={`h-1.5 w-1.5 rounded-full ${cfg.dot}`} />
      {cfg.label}
    </span>
  )
}

// ─── Spinner ───────────────────────────────────────────────────────────────────

function Spinner({ size = 14 }: { size?: number }) {
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

// ─── Action Button ─────────────────────────────────────────────────────────────

interface ActionBtnProps {
  label: string
  onClick: () => void
  loading?: boolean
  variant: 'blue' | 'green' | 'gray' | 'red'
}

const VARIANT_CLASS: Record<string, string> = {
  blue:  'bg-indigo-500 hover:bg-indigo-600 text-white',
  green: 'bg-emerald-500 hover:bg-emerald-600 text-white',
  gray:  'bg-slate-400 hover:bg-slate-500 text-white',
  red:   'bg-red-500 hover:bg-red-600 text-white',
}

function ActionBtn({ label, onClick, loading = false, variant }: ActionBtnProps) {
  return (
    <button
      onClick={onClick}
      disabled={loading}
      className={`inline-flex items-center gap-1 rounded-lg px-2.5 py-1 text-xs font-semibold shadow-sm
        transition-all duration-150 active:scale-95 disabled:opacity-60 disabled:cursor-not-allowed
        ${VARIANT_CLASS[variant]}`}
    >
      {loading && <Spinner size={11} />}
      {label}
    </button>
  )
}

// ─── Main Page ─────────────────────────────────────────────────────────────────

export default function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [queues, setQueues] = useState<Queue[]>([])
  const [estimatesMap, setEstimatesMap] = useState<Map<string, QueueETA>>(new Map())
  const [avgServiceMin, setAvgServiceMin] = useState<number>(30)
  const [date, setDate] = useState(format(new Date(), 'yyyy-MM-dd'))
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)
  // Track which rows are being updated (by queue id)
  const [updatingIds, setUpdatingIds] = useState<Set<string>>(new Set())
  const wsRef = useRef<WebSocket | null>(null)

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const [statsData, queuesData, etaData] = await Promise.all([
        getDashboardStats(date),
        getQueues(date, page, 10),
        getEstimates(date).catch(() => null),
      ])
      setStats(statsData)
      setQueues(queuesData.queues)
      setTotalPages(queuesData.total_pages)

      if (etaData) {
        setAvgServiceMin(etaData.avg_service_minutes || 30)
        const map = new Map<string, QueueETA>()
        etaData.estimates?.forEach((item: QueueETA) => map.set(item.id, item))
        setEstimatesMap(map)
      }
    } finally {
      setLoading(false)
    }
  }, [date, page])

  // WebSocket — listen for realtime events
  useEffect(() => {
    const ws = new WebSocket(
      process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws'
    )
    wsRef.current = ws
    ws.onmessage = (e) => {
      const event = JSON.parse(e.data)
      if (event.event === 'queue_created' || event.event === 'status_updated') {
        fetchData()
      }
    }
    ws.onclose = () => console.log('WebSocket disconnected')
    return () => ws.close()
  }, [fetchData])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  // ── Optimistic status update ─────────────────────────────────────────────────
  const handleStatusUpdate = async (id: string, status: QueueStatus) => {
    // 1. Mark as loading
    setUpdatingIds(prev => new Set(prev).add(id))

    // 2. Optimistically update local state immediately
    const previousQueues = queues
    setQueues(prev =>
      prev.map(q => q.id === id ? { ...q, status } : q)
    )

    // 3. Optimistically update stats
    // (just refresh stats silently in bg — no full reload)
    try {
      await updateQueueStatus(id, status)
      // On success: silently refresh stats & estimates
      getDashboardStats(date).then(setStats).catch(() => {})
      getEstimates(date).then(res => {
        if (res) {
          setAvgServiceMin(res.avg_service_minutes || 30)
          const map = new Map<string, QueueETA>()
          res.estimates?.forEach((item: QueueETA) => map.set(item.id, item))
          setEstimatesMap(map)
        }
      }).catch(() => {})
    } catch (err) {
      // On failure: revert local state
      console.error('Status update failed', err)
      setQueues(previousQueues)
    } finally {
      setUpdatingIds(prev => {
        const next = new Set(prev)
        next.delete(id)
        return next
      })
    }
  }

  // ── Optimistic delete ────────────────────────────────────────────────────────
  const handleDelete = async (id: string) => {
    if (!confirm('Yakin hapus antrian ini?')) return

    setUpdatingIds(prev => new Set(prev).add(id))
    const previousQueues = queues
    // Optimistically remove from list
    setQueues(prev => prev.filter(q => q.id !== id))

    try {
      await deleteQueue(id)
      getDashboardStats(date).then(setStats).catch(() => {})
    } catch (err) {
      console.error('Delete failed', err)
      setQueues(previousQueues)
    } finally {
      setUpdatingIds(prev => {
        const next = new Set(prev)
        next.delete(id)
        return next
      })
    }
  }

  const filteredQueues = queues.filter(q =>
    q.vehicle_plate.toLowerCase().includes(search.toLowerCase()) ||
    q.owner_name.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <main className="min-h-screen bg-slate-50 p-6">
      <div className="max-w-6xl mx-auto space-y-6">

        {/* ── Header ─────────────────────────────────────────────────────── */}
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 className="text-2xl font-extrabold text-slate-800 tracking-tight">
              🚗 Dashboard Antrian
            </h1>
            <p className="text-sm text-slate-500 mt-0.5">Showroom Queue Management System</p>
          </div>
          <div className="flex items-center gap-3">
            <span className="inline-flex items-center gap-1.5 rounded-xl border border-indigo-100 bg-indigo-50/80 px-3 py-1.5 text-xs font-semibold text-indigo-700 shadow-sm">
              ✨ Smart ETA: Rata-rata {avgServiceMin} mnt / antrian
            </span>
            <input
              type="date"
              value={date}
              onChange={e => { setDate(e.target.value); setPage(1) }}
              className="rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 shadow-sm
                focus:outline-none focus:ring-2 focus:ring-indigo-400 transition"
            />
          </div>
        </div>

        {/* ── Stats Cards ─────────────────────────────────────────────────── */}
        {stats && (
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
            <StatsCard label="Total"    value={stats.total}     color="bg-gradient-to-br from-slate-700 to-slate-900"       icon="📋" />
            <StatsCard label="Menunggu" value={stats.waiting}   color="bg-gradient-to-br from-amber-400 to-orange-500"      icon="⏳" />
            <StatsCard label="Diproses" value={stats.processed} color="bg-gradient-to-br from-indigo-500 to-indigo-700"     icon="⚙️" />
            <StatsCard label="Selesai"  value={stats.done}      color="bg-gradient-to-br from-emerald-400 to-emerald-600"   icon="✅" />
            <StatsCard label="Batal"    value={stats.cancelled}  color="bg-gradient-to-br from-red-400 to-red-600"           icon="❌" />
          </div>
        )}

        {/* ── AI Insight ───────────────────────────────────────────────────── */}
        <AISummaryCard date={date} />

        {/* ── Toolbar ──────────────────────────────────────────────────────── */}
        <div className="flex gap-3 items-center">
          <div className="relative flex-1">
            <span className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 text-sm">🔍</span>
            <input
              type="text"
              placeholder="Cari plat atau nama pemilik..."
              value={search}
              onChange={e => setSearch(e.target.value)}
              className="w-full rounded-xl border border-slate-200 bg-white pl-9 pr-4 py-2.5 text-sm shadow-sm
                focus:outline-none focus:ring-2 focus:ring-indigo-400 transition"
            />
          </div>
          <button
            onClick={fetchData}
            className="rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-sm font-medium text-slate-600
              hover:bg-slate-50 shadow-sm transition active:scale-95"
          >
            🔄 Refresh
          </button>
        </div>

        {/* ── Queue Table ──────────────────────────────────────────────────── */}
        <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-100 bg-slate-50/80">
                <th className="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">No</th>
                <th className="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">Plat</th>
                <th className="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">Nama</th>
                <th className="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">Status</th>
                <th className="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">Waktu Masuk</th>
                <th className="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">Est. Selesai</th>
                <th className="px-5 py-3.5 text-left text-xs font-semibold uppercase tracking-wider text-slate-500">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {loading ? (
                // Skeleton rows
                Array.from({ length: 5 }).map((_, i) => (
                  <tr key={i}>
                    {Array.from({ length: 7 }).map((__, j) => (
                      <td key={j} className="px-5 py-3.5">
                        <div className="h-4 rounded-md bg-slate-100 animate-pulse" style={{ width: `${60 + j * 5}%` }} />
                      </td>
                    ))}
                  </tr>
                ))
              ) : filteredQueues.length === 0 ? (
                <tr>
                  <td colSpan={7} className="py-16 text-center">
                    <div className="flex flex-col items-center gap-2 text-slate-400">
                      <span className="text-4xl">📭</span>
                      <p className="text-sm font-medium">Tidak ada data antrian</p>
                    </div>
                  </td>
                </tr>
              ) : filteredQueues.map(q => {
                const isUpdating = updatingIds.has(q.id)
                const etaInfo = estimatesMap.get(q.id)

                let etaDisplay = '-'
                let isOverdue = false
                if (etaInfo && (q.status === 'waiting' || q.status === 'processing')) {
                  const doneTime = new Date(etaInfo.estimated_done_at)
                  etaDisplay = `± ${format(doneTime, 'HH:mm')}`
                  isOverdue = isBefore(doneTime, new Date())
                }

                return (
                  <tr
                    key={q.id}
                    className={`group row-enter transition-all duration-200 hover:bg-indigo-50/40
                      ${isUpdating ? 'opacity-60 pointer-events-none' : ''}`}
                  >
                    {/* Queue number */}
                    <td className="px-5 py-3.5">
                      <span className="inline-flex items-center justify-center rounded-lg bg-indigo-50 px-2.5 py-1
                        text-xs font-bold text-indigo-600 ring-1 ring-indigo-100">
                        {q.queue_number}
                      </span>
                    </td>

                    {/* Plate */}
                    <td className="px-5 py-3.5">
                      <span className="font-bold tracking-widest text-slate-800">{q.vehicle_plate}</span>
                    </td>

                    {/* Owner */}
                    <td className="px-5 py-3.5 text-slate-600">{q.owner_name || '-'}</td>

                    {/* Status — transition-all ensures smooth color swap */}
                    <td className="px-5 py-3.5">
                      {isUpdating ? (
                        <span className="inline-flex items-center gap-1.5 text-slate-400 text-xs">
                          <Spinner size={13} /> Memperbarui...
                        </span>
                      ) : (
                        <StatusBadge status={q.status} />
                      )}
                    </td>

                    {/* Waktu Masuk */}
                    <td className="px-5 py-3.5 text-slate-500 tabular-nums">
                      {format(new Date(q.created_at), 'HH:mm')}
                    </td>

                    {/* Est. Selesai (AI Feature) */}
                    <td className="px-5 py-3.5 tabular-nums">
                      {q.status === 'done' ? (
                        <span className="text-emerald-600 font-medium text-xs">✓ Selesai</span>
                      ) : q.status === 'cancelled' ? (
                        <span className="text-slate-400 text-xs">-</span>
                      ) : (
                        <span className={`inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-xs font-semibold ${
                          isOverdue
                            ? 'bg-rose-50 text-rose-600 ring-1 ring-rose-200'
                            : 'bg-indigo-50 text-indigo-700 ring-1 ring-indigo-100'
                        }`}>
                          ⏳ {etaDisplay}
                        </span>
                      )}
                    </td>

                    {/* Actions */}
                    <td className="px-5 py-3.5">
                      <div className="flex gap-1.5">
                        {q.status === 'waiting' && (
                          <ActionBtn
                            label="Proses"
                            variant="blue"
                            loading={isUpdating}
                            onClick={() => handleStatusUpdate(q.id, 'processing')}
                          />
                        )}
                        {q.status === 'processing' && (
                          <ActionBtn
                            label="Selesai"
                            variant="green"
                            loading={isUpdating}
                            onClick={() => handleStatusUpdate(q.id, 'done')}
                          />
                        )}
                        {(q.status === 'waiting' || q.status === 'processing') && (
                          <ActionBtn
                            label="Batal"
                            variant="gray"
                            loading={isUpdating}
                            onClick={() => handleStatusUpdate(q.id, 'cancelled')}
                          />
                        )}
                        {(q.status === 'waiting' || q.status === 'cancelled') && (
                          <ActionBtn
                            label="Hapus"
                            variant="red"
                            loading={isUpdating}
                            onClick={() => handleDelete(q.id)}
                          />
                        )}
                      </div>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>

        {/* ── Pagination ───────────────────────────────────────────────────── */}
        <div className="flex items-center justify-between text-sm text-slate-500">
          <span>Halaman <strong className="text-slate-700">{page}</strong> dari <strong className="text-slate-700">{totalPages}</strong></span>
          <div className="flex gap-2">
            <button
              onClick={() => setPage(p => Math.max(1, p - 1))}
              disabled={page === 1}
              className="rounded-xl border border-slate-200 bg-white px-4 py-1.5 font-medium shadow-sm
                hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed transition active:scale-95"
            >
              ← Prev
            </button>
            <button
              onClick={() => setPage(p => Math.min(totalPages, p + 1))}
              disabled={page === totalPages}
              className="rounded-xl border border-slate-200 bg-white px-4 py-1.5 font-medium shadow-sm
                hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed transition active:scale-95"
            >
              Next →
            </button>
          </div>
        </div>

      </div>
    </main>
  )
}