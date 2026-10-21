import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// Сборка кладётся в ../server/dist: Go встраивает её через //go:embed
// и раздаёт из одного бинарника.
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    // В разработке фронтенд живёт отдельно, поэтому /api проксируется на Go.
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: '../server/dist',
    emptyOutDir: true,
    target: 'es2020',
  },
})
