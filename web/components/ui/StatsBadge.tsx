import { QueueStatus } from '@/lib/types'

const config: Record<QueueStatus, { label: string; className: string }> = {
  waiting:    { label: 'Menunggu',  className: 'bg-yellow-100 text-yellow-800' },
  processing: { label: 'Diproses', className: 'bg-blue-100 text-blue-800' },
  done:       { label: 'Selesai',  className: 'bg-green-100 text-green-800' },
  cancelled:  { label: 'Batal',    className: 'bg-red-100 text-red-800' },
}

export default function StatusBadge({ status }: { status: QueueStatus }) {
  const { label, className } = config[status]
  return (
    <span className={`px-2 py-1 rounded-full text-xs font-medium ${className}`}>
      {label}
    </span>
  )
}