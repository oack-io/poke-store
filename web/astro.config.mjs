import { defineConfig } from 'astro/config';

export default defineConfig({
  outDir: './dist',
  build: {
    assets: '_assets',
  },
});
