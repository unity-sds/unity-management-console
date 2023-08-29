import { sveltekit } from '@sveltejs/kit/vite';
import type { UserConfig } from 'vite';
import { defineConfig } from 'vite';

const config: UserConfig = {
        plugins: [sveltekit()],
};

export default defineConfig({
  ...config,
  base: '/ui/',
});
