import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd())

  return {
    envDir: './',
    plugins: [vue(), vueDevTools(), tailwindcss()],
    server: {
      proxy: {
        '/nodes': {
          changeOrigin: true,
          target: 'http://localhost:9095',
        },
        '/pipeline/start': {
          changeOrigin: true,
          target: 'http://localhost:9095',
        },
        '/updates/pipeline': {
          changeOrigin: true,
          target: 'ws://localhost:9095',
        },
      },
      fs: {
        allow: [`${env.VITE_SOURCE_DIR}godel`, `${env.VITE_SOURCE_DIR}Vue98`, './src'],
      },
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
  }
})
