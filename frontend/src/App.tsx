import { useCallback, useEffect, useState } from 'react'
import { ApiError, getDbHealth, getHealth } from './api/client'

type Status = 'loading' | 'up' | 'down'

function useStatus(fetcher: (s: AbortSignal) => Promise<unknown>) {
  const [status, setStatus] = useState<Status>('loading')
  const [detail, setDetail] = useState('')

  const load = useCallback(() => {
    const ctrl = new AbortController()
    fetcher(ctrl.signal)
      .then((r) => {
        setStatus('up')
        setDetail(JSON.stringify(r))
      })
      .catch((e: unknown) => {
        if (ctrl.signal.aborted) return
        setStatus('down')
        setDetail(e instanceof ApiError ? `${e.status} ${e.code}: ${e.message}` : String(e))
      })
    return () => ctrl.abort()
  }, [fetcher])

  useEffect(() => load(), [load])

  const reload = () => {
    setStatus('loading')
    load()
  }
  return { status, detail, reload }
}

function StatusRow({ label, status, detail }: { label: string; status: Status; detail: string }) {
  return (
    <div className="row">
      <span className={`dot dot-${status}`} aria-hidden />
      <div>
        <div className="label">{label}</div>
        <code className="detail">{status === 'loading' ? 'checking…' : detail}</code>
      </div>
    </div>
  )
}

export default function App() {
  const api = useStatus(getHealth)
  const db = useStatus(getDbHealth)

  return (
    <main className="container">
      <h1>Backend status</h1>
      <p className="muted">
        React (Vite) → <code>/health</code> via dev proxy → Go API → PostgreSQL
      </p>

      <section className="card">
        <StatusRow label="API  GET /health" {...api} />
        <StatusRow label="Database  GET /health/db" {...db} />
        <button
          onClick={() => {
            api.reload()
            db.reload()
          }}
        >
          Re-check
        </button>
      </section>

      {/* Feature UI goes here, e.g. <ItemsPage /> using api<T>() from ./api/client */}
    </main>
  )
}
