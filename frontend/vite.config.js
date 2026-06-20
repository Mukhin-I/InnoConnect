import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    historyApiFallback: true,
    host: '0.0.0.0',
    port: 80, // или 8080, если без sudo
  },
})
