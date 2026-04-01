import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import ListSources from './ListSources.svelte';

const { fetchPaginatedData } = vi.hoisted(() => ({
	fetchPaginatedData: vi.fn()
}));

vi.mock('$app/paths', () => ({
	resolve: (path: string) => path
}));

vi.mock('$lib/api', () => ({
	fetchPaginatedData
}));

describe('ListSources', () => {
	it('loads and renders sources on mount', async () => {
		fetchPaginatedData.mockReset();
		fetchPaginatedData.mockResolvedValueOnce({
			data: [
				{
					id: 7,
					public_id: 'src_7',
					ingress_url: 'http://localhost:3000/api/v1/ingest/src_7',
					egress_url: 'https://example.com/out',
					static_headers: { Authorization: 'Bearer token' },
					status: 'active',
					status_reason: 'ready',
					description: 'Primary source',
					created_at: '2026-04-01T00:00:00.000Z',
					updated_at: '2026-04-02T00:00:00.000Z',
					disable_at: null
				}
			],
			next_cursor: null,
			has_next: false
		});

		render(ListSources);

		expect(await screen.findByText('Primary source')).toBeInTheDocument();
		expect(screen.getByRole('link', { name: '7' })).toHaveAttribute('href', '/sources/7');
		expect(fetchPaginatedData).toHaveBeenCalledWith(
			'/api/sources',
			20,
			null,
			expect.objectContaining({
				filter_status: '*',
				sort_direction: 'DESC'
			})
		);
	});

	it('requests another page when load more is clicked', async () => {
		fetchPaginatedData.mockReset();
		fetchPaginatedData
			.mockResolvedValueOnce({
				data: [
					{
						id: 7,
						public_id: 'src_7',
						ingress_url: 'http://localhost:3000/api/v1/ingest/src_7',
						egress_url: 'https://example.com/out',
						static_headers: {},
						status: 'active',
						created_at: '2026-04-01T00:00:00.000Z',
						updated_at: '2026-04-02T00:00:00.000Z'
					}
				],
				next_cursor: 'next-cursor',
				has_next: true
			})
			.mockResolvedValueOnce({
				data: [
					{
						id: 8,
						public_id: 'src_8',
						ingress_url: 'http://localhost:3000/api/v1/ingest/src_8',
						egress_url: 'https://example.com/out-2',
						static_headers: {},
						status: 'paused',
						created_at: '2026-04-03T00:00:00.000Z',
						updated_at: '2026-04-04T00:00:00.000Z'
					}
				],
				next_cursor: null,
				has_next: false
			});

		render(ListSources);

		await screen.findByText('7');
		await fireEvent.click(screen.getByRole('button', { name: 'Load More' }));

		expect(await screen.findByText('8')).toBeInTheDocument();
		expect(fetchPaginatedData).toHaveBeenNthCalledWith(
			2,
			'/api/sources',
			20,
			'next-cursor',
			expect.objectContaining({
				filter_status: '*',
				sort_direction: 'DESC'
			})
		);
	});
});
