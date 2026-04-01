import { afterEach, describe, expect, it, vi } from 'vitest';
import { fetchPaginatedData } from './api';

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
				statusText: 'Bad Request'
			})
		);

		await expect(fetchPaginatedData('/api/sources', 20, null)).rejects.toThrow(
			'Failed to fetch data: Bad Request'
		);
	});
});
