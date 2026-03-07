import { defineConfig } from "vite";

export default defineConfig({
  build: {
    outDir: "web/assets/js/dist",
    emptyOutDir: true,
    rollupOptions: {
      input: {
        tasks: "web/assets/js/apps/tasks/tasks.entry.js",
      },
      output: {
        entryFileNames: "[name].min.js",
      },
    },
  },
});
