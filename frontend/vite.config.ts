import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// Dev server proxies API calls to the Go backend, so the browser talks to one
// origin (no CORS issues locally). Override the target with VITE_API_TARGET.
const target = process.env.VITE_API_TARGET ?? 'http://localhost:8080'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': { target, changeOrigin: true },
      '/health': { target, changeOrigin: true },
    },
  },
})
