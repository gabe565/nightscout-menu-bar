import vue from "@vitejs/plugin-vue";
import { URL, fileURLToPath } from "node:url";
import Icons from "unplugin-icons/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [vue(), Icons({ compiler: "vue3" })],
  server: {
    // Wails proxies to the dev server over IPv4, and the page's origin is wails://
    host: "127.0.0.1",
    port: 5173,
    strictPort: true,
    hmr: { host: "127.0.0.1" },
  },
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
});
