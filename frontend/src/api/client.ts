// Thin fetch wrapper matching the backend's JSON contract:
//   success → the JSON body
//   error   → {"error":{"code","message","details?"}}  (thrown as ApiError)

export type FieldError = { field: string; issue: string }

export class ApiError extends Error {
  status: number
  code: string
  details?: FieldError[]

  constructor(status: number, code: string, message: string, details?: FieldError[]) {
    super(message)
    this.status = status
    this.code = code
    this.details = details
  }
}

// next_cursor is "" on the last page; pass it back as ?cursor= for the next one.
export type Page<T> = { data: T[]; next_cursor: string; limit: number }

type Options = {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  body?: unknown
  userId?: string // sent as X-User-ID when the API identifies callers by header
  query?: Record<string, string | number | undefined>
  signal?: AbortSignal
}

const BASE = import.meta.env.VITE_API_BASE ?? '' // '' → same origin via Vite proxy

export async function api<T>(path: string, opts: Options = {}): Promise<T> {
  const url = new URL(BASE + path, window.location.origin)
  for (const [k, v] of Object.entries(opts.query ?? {})) {
    if (v !== undefined && v !== '') url.searchParams.set(k, String(v))
  }

  const headers: Record<string, string> = { Accept: 'application/json' }
  if (opts.body !== undefined) headers['Content-Type'] = 'application/json'
  if (opts.userId) headers['X-User-ID'] = opts.userId

  const res = await fetch(url, {
    method: opts.method ?? 'GET',
    headers,
    body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
    signal: opts.signal,
  })

  if (res.status === 204) return undefined as T

  const text = await res.text()
  const json = text ? JSON.parse(text) : undefined

  if (!res.ok) {
    const e = json?.error
    throw new ApiError(res.status, e?.code ?? 'UNKNOWN', e?.message ?? res.statusText, e?.details)
  }
  return json as T
}

// ---- endpoints --------------------------------------------------------------

export type Health = { status: string; database?: string }

export const getHealth = (signal?: AbortSignal) => api<Health>('/health', { signal })
export const getDbHealth = (signal?: AbortSignal) => api<Health>('/health/db', { signal })
