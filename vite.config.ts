import { sveltekit } from '@sveltejs/kit/vite';
import type { UserConfig } from 'vite';
import { defineConfig } from 'vite';
import purgeCss from 'vite-plugin-svelte-purgecss';

const config: UserConfig = {
	plugins: [sveltekit(), purgeCss()],
	css: {
		postcss: './postcss.config.cjs'
	},
	ssr: {
		noExternal: ['xterm', 'xterm-addon-fit']
	}
};

export default defineConfig({
	...config,
	base: '/ui/'
});
