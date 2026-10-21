import type { Content } from '../types'

export interface AdminDocument {
  ru: Content
  en: Content
}

export interface AdminSnapshot {
  content: AdminDocument
  rev: number
  savedAt: string
}

export class AdminError extends Error {
  status: number
  constructor(message: string, status: number) {
    super(message)
    this.name = 'AdminError'
    this.status = status
  }
}

async function request<T>(url: string, init: RequestInit = {}): Promise<T> {
  const res = await fetch(url, {
    credentials: 'same-origin',
    ...init,
    headers: {
      Accept: 'application/json',
      ...(init.body ? { 'Content-Type': 'application/json' } : {}),
      // Заголовок обязателен для записи: он же служит CSRF-защитой.
      ...(init.method && init.method !== 'GET' ? { 'X-Portfolio-Admin': '1' } : {}),
      ...(init.headers ?? {}),
    },
  })

  if (res.status === 204) return undefined as T

  const text = await res.text()
  let payload: unknown = null
  try {
    payload = text ? JSON.parse(text) : null
  } catch {
    payload = null
  }

  if (!res.ok) {
    const message =
      (payload && typeof payload === 'object' && 'error' in payload
        ? String((payload as { error: unknown }).error)
        : '') ||
      text ||
      `Запрос ${url} вернул ${res.status}`
    throw new AdminError(message, res.status)
  }

  return payload as T
}

export const adminApi = {
  session: () => request<{ authenticated: boolean }>('/api/admin/session'),

  login: (password: string) =>
    request<{ ok: boolean }>('/api/admin/login', {
      method: 'POST',
      body: JSON.stringify({ password }),
    }),

  logout: () => request<{ ok: boolean }>('/api/admin/logout', { method: 'POST' }),

  load: () => request<AdminSnapshot>('/api/admin/content'),

  save: (doc: AdminDocument) =>
    request<{ ok: boolean; rev: number; savedAt: string }>('/api/admin/content', {
      method: 'PUT',
      body: JSON.stringify(doc),
    }),
}
