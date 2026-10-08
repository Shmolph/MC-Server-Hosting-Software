import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

const outputDirectory = process.env.VITE_OUTPUT_DIR || "dist";

export default defineConfig({
  plugins: [svelte()],
  server: {
    host: "localhost",
    port: 5000,
    strictPort: true
  },
  build: {
    outDir: outputDirectory,
    emptyOutDir: true,
    rollupOptions: {
      input: "index.html"
    }
  }
});
