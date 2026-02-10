import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import env from './env';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: `http://localhost:${env.API_PORT}/api/v1`,
				changeOrigin: true,
				secure: false,
				rewrite: (path) => path.replace(/^\/api/, '')
			}
		}
	}
});
