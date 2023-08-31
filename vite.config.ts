import { sveltekit } from '@sveltejs/kit/vite';
import type { UserConfig } from 'vite';
import { defineConfig } from 'vite';

const config: UserConfig = {
	plugins: [sveltekit()],
	ssr: {
		noExternal: ['xterm', 'xterm-addon-fit']
	}
};

export default defineConfig({
	...config,
	base: '/ui/'
});
