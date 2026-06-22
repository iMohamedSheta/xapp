import { defineConfig } from 'vite';
import laravel from 'laravel-vite-plugin';
import vue from '@vitejs/plugin-vue';
import path from 'path';
import { VitePWA } from 'vite-plugin-pwa';

export default defineConfig({
  base: './',
  plugins: [
    laravel({
      input: 'resources/js/app.ts',
      publicDirectory: 'public',
      buildDirectory: 'build',
      refresh: true,
      ssr: 'resources/js/ssr.ts',
    }),
    vue({
      include: [/\.vue$/],
    }),
    // --- PWA Plugin ---
    VitePWA({
      registerType: 'autoUpdate',
      manifest: {
        "name": "XApp",
        "short_name": "XApp",
        "start_url": "/",
        "theme_color": "#1f1f1f",
        "background_color": "#1f1f1f",
        "display": "standalone",
        "lang": "en",
        "scope": "./",
        "icons": [
          {
            "src": "/public/favicon/web-app-manifest-96x96.png",
            "sizes": "96x96",
            "type": "image/png"
          },
          {
            "src": "/public/favicon/web-app-manifest-192x192.png",
            "sizes": "192x192",
            "type": "image/png"
          },
          {
            "src": "/public/favicon/web-app-manifest-512x512.png",
            "sizes": "512x512",
            "type": "image/png"
          }
        ]
      },
      workbox: {
        globPatterns: ['**/*.{js,css,html,png,svg,ico}'],
      },
    })
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'resources/js'),
    },
  },
  build: {
    manifest: true, // Generate manifest.json file
    // outDir: 'public/build',
    rollupOptions: {
      input: 'resources/js/app.ts',
      output: {
        entryFileNames: 'assets/[name].js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name].[ext]',
        manualChunks: undefined, // Disable automatic chunk splitting
      },
    },
  },
  server: {
    host: '0.0.0.0', // bind to all interfaces
    port: 5173,
    strictPort: true,
    hmr: {
      host: '192.168.1.38', // your main LAN IP
      port: 5173,
    },
    cors: {
      origin: [
        'http://localhost:8080',
        'http://127.0.0.1:8080'
      ],
      credentials: true,
    },
    watch: {
      usePolling: true, // Set to true if you're on Windows/WSL
      ignored: [
        '**/node_modules/**',
        '**/dist/**',
        '**/storage/**',        // Ignore storage directory
        '**/public/build/**',   // Ignore build directory
        '**/tmp/**',            // Ignore temp directories
        '**/*.tmp',             // Ignore temp files
        '**/*.log'              // Ignore log files
      ],
    },
  },
});
