import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  return {
    plugins: [react()],
    build: {
      outDir: mode === "docker" ? "static" : "../static",
      emptyOutDir: true, // also necessary
    },
    resolve: {
      alias: {
        components: "/src/Components"
      }
    }
  }
})

