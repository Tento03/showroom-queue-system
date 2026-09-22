interface StatsCardProps {
  label: string
  value: number
  color: string
  icon: string
  trend?: string
}

export default function StatsCard({ label, value, color, icon, trend }: StatsCardProps) {
  return (
    <div className={`relative overflow-hidden rounded-2xl p-5 text-white shadow-lg ${color} transition-transform duration-200 hover:-translate-y-0.5 hover:shadow-xl`}>
      {/* Background decoration */}
      <div className="absolute -right-4 -top-4 text-6xl opacity-10 select-none">{icon}</div>

      <p className="text-xs font-semibold uppercase tracking-wider opacity-75">{label}</p>
      <p className="text-4xl font-extrabold mt-1 tabular-nums">{value}</p>
      {trend && (
        <p className="text-xs mt-2 opacity-80">{trend}</p>
      )}
    </div>
  )
}