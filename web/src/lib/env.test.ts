import { afterEach, describe, expect, it, vi } from 'vitest';

describe('env', () => {
	afterEach(() => {
		vi.unstubAllEnvs();
		vi.resetModules();
	});

	it('uses provided frontend env values', async () => {
		vi.stubEnv('VITE_ENV', 'uat');
		vi.stubEnv('VITE_API_BASE_URL', 'https://api.example.com/api/v1');

		const env = (await import('./env')).default;

		expect(env).toEqual({
			VITE_ENV: 'uat',
			VITE_API_BASE_URL: 'https://api.example.com/api/v1'
		});
	});

	it('falls back to defaults when frontend env values are unset', async () => {
		vi.stubEnv('VITE_ENV', undefined);
		vi.stubEnv('VITE_API_BASE_URL', undefined);

		const env = (await import('./env')).default;

		expect(env).toEqual({
			VITE_ENV: 'dev',
			VITE_API_BASE_URL: 'http://localhost:3000/api/v1'
		});
	});
});
