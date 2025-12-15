import axios from 'axios'
import type { Resource, ResourcePayload, ReviewPayload, UserProfile } from '@/types'

const api = axios.create({
  baseURL: '/api',
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export async function fetchResources(params: Record<string, unknown> = {}): Promise<Resource[]> {
  const { data } = await api.get<Resource[]>('/resources', { params })
  return data
}

export async function fetchResource(id: string): Promise<Resource> {
  const { data } = await api.get<Resource>(`/resources/${id}`)
  return data
}

export async function fetchRecommendations(id: string): Promise<Resource[]> {
  const { data } = await api.get<Resource[]>(`/resources/${id}/recommendations`)
  return data
}

export async function fetchVersions(id: string): Promise<Resource[]> {
  const { data } = await api.get<Resource[]>(`/resources/${id}/versions`)
  return data
}

export async function uploadResource(payload: ResourcePayload): Promise<Resource> {
  const form = new FormData()
  Object.entries(payload).forEach(([key, value]) => {
    if (value !== undefined && value !== null) form.append(key, value as Blob | string)
  })
  const { data } = await api.post<Resource>('/resources', form, { headers: { 'Content-Type': 'multipart/form-data' } })
  return data
}

export async function downloadResource(id: string): Promise<void> {
  const { data, headers } = await api.get(`/resources/${id}/download`, { responseType: 'blob' })
  const blob = new Blob([data])
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = headers['content-disposition']?.split('filename=')[1] || 'resource'
  a.click()
  window.URL.revokeObjectURL(url)
}

export async function submitReview(id: string, payload: ReviewPayload) {
  await api.post(`/resources/${id}/reviews`, payload)
}

export async function login(email: string, password: string): Promise<{ token: string; user: UserProfile }> {
  const { data } = await api.post('/auth/login', { email, password })
  return data
}

export async function register(
  email: string,
  password: string,
  displayName: string,
): Promise<{ token: string; user: UserProfile }> {
  const { data } = await api.post('/auth/register', { email, password, displayName })
  return data
}

export async function toggleFavorite(id: string): Promise<{ status: string }> {
  const { data } = await api.post(`/resources/${id}/favorite`)
  return data
}

export async function fetchFavorites(): Promise<Resource[]> {
  const { data } = await api.get<Resource[]>('/user/favorites')
  return data
}

export async function fetchDownloads(): Promise<Resource[]> {
  const { data } = await api.get<Resource[]>('/user/downloads')
  return data
}

export async function reportResource(id: string, reason: string) {
  await api.post(`/resources/${id}/report`, { reason })
}

export async function fetchPendingResources(): Promise<Resource[]> {
  const { data } = await api.get<Resource[]>('/admin/pending')
  return data
}

export async function auditResource(id: string, action: 'approve' | 'reject', reason?: string) {
  await api.post(`/admin/resources/${id}/audit`, { action, reason })
}

export async function changePassword(oldPassword: string, newPassword: string) {
  await api.post('/user/change-password', { oldPassword, newPassword })
}

export async function fetchReports(): Promise<any[]> {
  const { data } = await api.get('/admin/reports')
  return data
}

export async function resolveReport(id: string) {
  await api.post(`/admin/reports/${id}/resolve`)
}

export async function fetchPopularTags(): Promise<{ tag: string; count: number }[]> {
  const { data } = await api.get<{ tag: string; count: number }[]>('/resources/tags/popular')
  return data
}

