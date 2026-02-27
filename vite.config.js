import { defineConfig } from "vite";

export default defineConfig({
  build: {
    outDir: "web/assets/css/dist",
    emptyOutDir: true,
    rollupOptions: {
      input: {
        tasks: "web/assets/css/apps/tasks.entry.css",
        home: "web/assets/css/public/home.entry.css",
      },
      output: {
        assetFileNames: "[name].min[extname]",
        entryFileNames: "_[name].js",
      },
    },
    cssMinify: true,
  },
});
