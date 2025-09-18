import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api/v1/orders': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
      '/api/v1/products': {
        target: 'http://localhost:8082',
        changeOrigin: true,
      },
      '/product-images': {
        target: 'http://localhost:9000',
        changeOrigin: true,
      },
    }
  }
})
