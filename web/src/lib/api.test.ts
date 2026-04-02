import { afterEach, describe, expect, it, vi } from 'vitest';
import { fetchPaginatedData, getResponseErrorMessage } from './api';

describe('fetchPaginatedData', () => {
	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('builds the request URL and normalizes the paginated response', async () => {
		const fetchMock = vi.fn().mockResolvedValue({
			ok: true,
			json: async () => ({
				data: [{ id: 1 }],
				next_cursor: 'next-token',
				has_next: true
			})
		});
		vi.stubGlobal('fetch', fetchMock);

		const result = await fetchPaginatedData<{ id: number }>('/api/events', 10, 'cursor-123', {
			search: 'invoice',
			state: 'failed'
		});

		expect(fetchMock).toHaveBeenCalledWith(
			'/api/events?limit=10&search=invoice&state=failed&cursor=cursor-123'
		);
		expect(result).toEqual({
			data: [{ id: 1 }],
			next_cursor: 'next-token',
			has_next: true
		});
	});

	it('defaults next_cursor to null when the API omits it', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({
					data: [],
					has_next: false
				})
			})
		);

		const result = await fetchPaginatedData('/api/sources', 20, null);

		expect(result).toEqual({
			data: [],
			next_cursor: null,
			has_next: false
		});
	});

	it('throws when the API response is not ok', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn().mockResolvedValue({
				ok: false,
				json: async () => ({ error: 'Sources request failed' }),
				statusText: 'Bad Request'
			})
		);

		await expect(fetchPaginatedData('/api/sources', 20, null)).rejects.toThrow(
			'Sources request failed'
		);
	});

	it('warns and falls back to statusText when the error response is not JSON', async () => {
		const consoleWarn = vi.spyOn(console, 'warn').mockImplementation(() => {});
		vi.stubGlobal(
			'fetch',
			vi.fn().mockResolvedValue({
				ok: false,
				status: 400,
				json: async () => {
					throw new Error('invalid json');
				},
				statusText: 'Bad Request'
			})
		);

		await expect(fetchPaginatedData('/api/sources', 20, null)).rejects.toThrow(
			'Failed to fetch data: 400 Bad Request'
		);
		expect(consoleWarn).toHaveBeenCalledWith(
			'Failed to parse error response as JSON',
			expect.objectContaining({ ok: false, status: 400, statusText: 'Bad Request' })
		);
	});

	it('includes the status code in the fallback error message when JSON has no error field', async () => {
		const response = {
			ok: false,
			status: 404,
			statusText: 'Not Found',
			json: async () => ({ detail: 'missing' })
		} as Response;

		await expect(getResponseErrorMessage(response, 'Failed to fetch source')).resolves.toBe(
			'Failed to fetch source: 404 Not Found'
		);
	});
});
