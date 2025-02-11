import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'

// https://vitejs.dev/config/
export default defineConfig({
  server :{
    host : '0.0.0.0'
  },
  plugins: [
    vue(),
    vueJsx(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  devServer:{
    host: '0.0.0.0',
    port: 9090,
    client: {
      webSocketURL: 'ws://localhost:9090/chat',
    },
    headers: {
      'Access-Control-Allow-Origin': '*',
    }
  }
})
