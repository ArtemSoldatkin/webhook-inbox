import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import env from './env';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: `${env.API_PROTOCOL}://${env.API_HOST}:${env.API_PORT}/api/v1`,
				changeOrigin: true,
				secure: false,
				rewrite: (path) => path.replace(/^\/api/, '')
			}
		}
	},
	define: {
		'import.meta.env.VITE_ENV': JSON.stringify(env.ENV)
	}
});
