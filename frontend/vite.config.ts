import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import tailwindcss from '@tailwindcss/vite';
import wails from '@wailsio/runtime/plugins/vite';
import { fileURLToPath } from 'url';

export default defineConfig({
  plugins: [vue(), tailwindcss(), wails('./bindings')],
  server: {
    host: '0.0.0.0',
    port: 9245,
    strictPort: true,
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      bindings: fileURLToPath(new URL('./bindings', import.meta.url)),
    },
  },
});
