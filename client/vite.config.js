import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    // builds to Nubayrah's root/static/ -> where goserver will host front-end
    outDir: '../static',

    emptyOutDir: true, // also necessary
  }
})

