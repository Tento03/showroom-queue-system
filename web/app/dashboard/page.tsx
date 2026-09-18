'use client'

import { useEffect, useState, useCallback } from 'react'
import { format } from 'date-fns'
import { Queue, DashboardStats } from '@/lib/types'
import { getQueues, getDashboardStats, updateQueueStatus, deleteQueue } from '@/lib/api'
import StatsCard from '@/components/ui/StatsCard'

function StatusBadge({ status }: { status: string }) {
  const labels: Record<string, string> = {
    waiting: 'Menunggu',
    processing: 'Diproses',
    done: 'Selesai',
    cancelled: 'Batal',
  }
  const colors: Record<string, string> = {
    waiting: 'bg-yellow-100 text-yellow-800',
    processing: 'bg-blue-100 text-blue-800',
    done: 'bg-green-100 text-green-800',
    cancelled: 'bg-red-100 text-red-800',
  }

  return (
    <span className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${colors[status] || 'bg-gray-100 text-gray-800'}`}>
      {labels[status] || status}
    </span>
  )
}

export default function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [queues, setQueues] = useState<Queue[]>([])
  const [date, setDate] = useState(format(new Date(), 'yyyy-MM-dd'))
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const [statsData, queuesData] = await Promise.all([
        getDashboardStats(date),
        getQueues(date, page, 10),
      ])
      setStats(statsData)
      setQueues(queuesData.queues)
      setTotalPages(queuesData.total_pages)
    } finally {
      setLoading(false)
    }
  }, [date, page])

  // WebSocket — listen untuk event realtime
  useEffect(() => {
    const ws = new WebSocket(
      process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws'
    )

    ws.onmessage = (e) => {
      const event = JSON.parse(e.data)
      if (event.event === 'queue_created' || event.event === 'status_updated') {
        fetchData() // refresh data tiap ada event
      }
    }

    ws.onclose = () => console.log('WebSocket disconnected')

    return () => ws.close()
  }, [fetchData])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  const handleStatusUpdate = async (id: string, status: string) => {
    await updateQueueStatus(id, status as any)
    fetchData()
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Yakin hapus antrian ini?')) return
    await deleteQueue(id)
    fetchData()
  }

  const filteredQueues = queues.filter(q =>
    q.vehicle_plate.toLowerCase().includes(search.toLowerCase()) ||
    q.owner_name.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <main className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-6xl mx-auto space-y-6">

        {/* Header */}
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold text-gray-800">
            Dashboard Antrian
          </h1>
          <input
            type="date"
            value={date}
            onChange={e => { setDate(e.target.value); setPage(1) }}
            className="border rounded-lg px-3 py-2 text-sm"
          />
        </div>

        {/* Stats Cards */}
        {stats && (
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
            <StatsCard label="Total"     value={stats.total}      color="bg-gray-700" />
            <StatsCard label="Menunggu"  value={stats.waiting}    color="bg-yellow-500" />
            <StatsCard label="Diproses" value={stats.processed} color="bg-blue-500" />
            <StatsCard label="Selesai"   value={stats.done}       color="bg-green-500" />
            <StatsCard label="Batal"     value={stats.cancelled}  color="bg-red-500" />
          </div>
        )}

        {/* Search */}
        <input
          type="text"
          placeholder="Cari plat atau nama..."
          value={search}
          onChange={e => setSearch(e.target.value)}
          className="w-full border rounded-lg px-4 py-2 text-sm"
        />

        {/* Queue Table */}
        <div className="bg-white rounded-xl shadow overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 text-gray-600">
              <tr>
                <th className="px-4 py-3 text-left">No</th>
                <th className="px-4 py-3 text-left">Plat</th>
                <th className="px-4 py-3 text-left">Nama</th>
                <th className="px-4 py-3 text-left">Status</th>
                <th className="px-4 py-3 text-left">Waktu</th>
                <th className="px-4 py-3 text-left">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {loading ? (
                <tr>
                  <td colSpan={6} className="text-center py-8 text-gray-400">
                    Loading...
                  </td>
                </tr>
              ) : filteredQueues.length === 0 ? (
                <tr>
                  <td colSpan={6} className="text-center py-8 text-gray-400">
                    Tidak ada data
                  </td>
                </tr>
              ) : filteredQueues.map(q => (
                <tr key={q.id} className="hover:bg-gray-50 transition-colors duration-200">
                  <td className="px-4 py-3 font-bold">{q.queue_number}</td>
                  <td className="px-4 py-3">{q.vehicle_plate}</td>
                  <td className="px-4 py-3">{q.owner_name || '-'}</td>
                  <td className="px-4 py-3">
                    <StatusBadge status={q.status} />
                  </td>
                  <td className="px-4 py-3 text-gray-500">
                    {format(new Date(q.created_at), 'HH:mm')}
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex gap-2">
                      {q.status === 'waiting' && (
                        <button
                          onClick={() => handleStatusUpdate(q.id, 'processing')}
                          className="text-xs bg-blue-500 hover:bg-blue-600 active:scale-95 transition-all duration-200 text-white px-2 py-1 rounded shadow hover:shadow-md"
                        >
                          Proses
                        </button>
                      )}
                      {q.status === 'processing' && (
                        <button
                          onClick={() => handleStatusUpdate(q.id, 'done')}
                          className="text-xs bg-green-500 hover:bg-green-600 active:scale-95 transition-all duration-200 text-white px-2 py-1 rounded shadow hover:shadow-md"
                        >
                          Selesai
                        </button>
                      )}
                      {(q.status === 'waiting' || q.status === 'processing') && (
                        <button
                          onClick={() => handleStatusUpdate(q.id, 'cancelled')}
                          className="text-xs bg-gray-400 hover:bg-gray-500 active:scale-95 transition-all duration-200 text-white px-2 py-1 rounded shadow hover:shadow-md"
                        >
                          Batal
                        </button>
                      )}
                      {(q.status === 'waiting' || q.status === 'cancelled') && (
                        <button
                          onClick={() => handleDelete(q.id)}
                          className="text-xs bg-red-500 hover:bg-red-600 active:scale-95 transition-all duration-200 text-white px-2 py-1 rounded shadow hover:shadow-md"
                        >
                          Hapus
                        </button>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div className="flex items-center justify-between text-sm text-gray-600">
          <span>Halaman {page} dari {totalPages}</span>
          <div className="flex gap-2">
            <button
              onClick={() => setPage(p => Math.max(1, p - 1))}
              disabled={page === 1}
              className="px-3 py-1 border rounded disabled:opacity-40"
            >
              ← Prev
            </button>
            <button
              onClick={() => setPage(p => Math.min(totalPages, p + 1))}
              disabled={page === totalPages}
              className="px-3 py-1 border rounded disabled:opacity-40"
            >
              Next →
            </button>
          </div>
        </div>

      </div>
    </main>
  )
}