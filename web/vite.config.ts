import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import fs from 'fs';
import { defineConfig } from 'vite';
import env from './env';

export default defineConfig({
	plugins: [sveltekit(), tailwindcss()],
	resolve: {
		conditions: process.env.VITEST ? ['browser'] : []
	},
	test: {
		environment: 'jsdom',
		setupFiles: ['./src/test/setup.ts'],
		include: ['src/**/*.{test,spec}.{ts,js}']
	},
	server: {
		host: env.UI_HOST,
		port: Number(env.UI_PORT),
		...(env.UI_PROTOCOL === 'https' && {
			https: {
				key: fs.readFileSync(env.UI_HTTPS_KEY_PATH),
				cert: fs.readFileSync(env.UI_HTTPS_CERT_PATH)
			}
		}),
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
		'import.meta.env.VITE_ENV': JSON.stringify(env.ENV),
		'import.meta.env.VITE_API_BASE_URL': JSON.stringify(
			`${env.API_PROTOCOL}://${env.API_HOST}:${env.API_PORT}/api/v1`
		)
	}
});
