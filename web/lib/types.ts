export type QueueStatus = 'waiting' | 'processing' | 'done' | 'cancelled'

export interface Queue {
  id: string
  queue_number: string
  queue_date: string
  vehicle_plate: string
  vehicle_image_url: string
  owner_name: string
  owner_phone: string
  status: QueueStatus
  created_at: string
}

export interface DashboardStats {
  total: number
  waiting: number
  processed: number
  done: number
  cancelled: number
}

export interface PaginatedQueues {
  date: string
  total: number
  page: number
  limit: number
  total_pages: number
  queues: Queue[]
}