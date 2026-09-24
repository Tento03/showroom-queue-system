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

export interface QueueETA {
  id: string
  queue_number: string
  vehicle_plate: string
  status: QueueStatus
  estimated_wait_minutes: number
  estimated_done_at: string
  queues_ahead: number
}

export interface ETAResponse {
  date: string
  avg_service_minutes: number
  estimates: QueueETA[]
}

export interface AISummaryResponse {
  summary: string
  generated_at: string
  cached: boolean
}
