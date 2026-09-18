interface StatsCardProps {
  label: string
  value: number
  color: string
}

export default function StatsCard({ label, value, color }: StatsCardProps) {
  return (
    <div className={`rounded-xl p-5 text-white ${color}`}>
      <p className="text-sm font-medium opacity-80">{label}</p>
      <p className="text-4xl font-bold mt-1">{value}</p>
    </div>
  )
}