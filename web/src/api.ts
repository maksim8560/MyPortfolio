import type { Content, Health } from './types'

async function getJSON<T>(url: string, signal?: AbortSignal): Promise<T> {
  const res = await fetch(url, {
    signal,
    headers: { Accept: 'application/json' },
  })
  if (!res.ok) {
    throw new Error(`Запрос ${url} вернул ${res.status}`)
  }
  return (await res.json()) as T
}

export function fetchContent(lang: string, signal?: AbortSignal): Promise<Content> {
  return getJSON<Content>(`/api/content?lang=${encodeURIComponent(lang)}`, signal)
}

export function fetchHealth(signal?: AbortSignal): Promise<Health> {
  return getJSON<Health>('/api/healthz', signal)
}
