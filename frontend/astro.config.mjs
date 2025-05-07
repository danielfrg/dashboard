// @ts-check
import alpinejs from '@astrojs/alpinejs';
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from 'astro/config';


// https://astro.build/config
export default defineConfig({
    integrations: [alpinejs()],
    vite: {
      plugins: [tailwindcss()],
    },
  });
