import axios from 'axios'
import { QueueStatus } from './types'

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
})

export const getQueues = async (date?: string, page = 1, limit = 10) => {
  const params = new URLSearchParams()
  if (date) params.append('date', date)
  params.append('page', page.toString())
  params.append('limit', limit.toString())
  const res = await api.get(`/queues?${params}`)
  return res.data
}

export const getDashboardStats = async (date?: string) => {
  const params = date ? `?date=${date}` : ''
  const res = await api.get(`/dashboard/stats${params}`)
  return res.data
}

export const updateQueueStatus = async (id: string, status: QueueStatus) => {
  const res = await api.patch(`/queue/${id}/status`, { status })
  return res.data
}

export const deleteQueue = async (id: string) => {
  const res = await api.delete(`/queue/${id}`)
  return res.data
}