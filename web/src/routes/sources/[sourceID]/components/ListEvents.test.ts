import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import ListEvents from './ListEvents.svelte';

const { fetchPaginatedData } = vi.hoisted(() => ({
	fetchPaginatedData: vi.fn()
}));

vi.mock('$app/paths', () => ({
	resolve: (path: string) => path
}));

vi.mock('$lib/api', () => ({
	fetchPaginatedData
}));

vi.mock('./BodyView.svelte', () => ({
	default: () => ({
		c: () => {},
		m: () => {},
		p: () => {},
		d: () => {}
	})
}));

describe('ListEvents', () => {
	it('loads and renders events on mount', async () => {
		fetchPaginatedData.mockReset();
		fetchPaginatedData.mockResolvedValueOnce({
			data: [
				{
					id: 11,
					source_id: 7,
					dedup_hash: 'abc123',
					method: 'POST',
					ingress_path: '/hook',
					query_params: { tenant: ['acme'] },
					raw_headers: { 'content-type': ['application/json'] },
					body: btoa('payload'),
					body_content_type: 'text/plain',
					received_at: '2026-04-01T00:00:00.000Z'
				}
			],
			next_cursor: null,
			has_next: false
		});

		render(ListEvents, {
			props: {
				sourceID: '7'
			}
		});

		expect(await screen.findByRole('link', { name: 'Event ID: 11' })).toHaveAttribute(
			'href',
			'/sources/7/11'
		);
		expect(screen.getByText('abc123')).toBeInTheDocument();
		expect(screen.getByText('tenant')).toBeInTheDocument();
		expect(fetchPaginatedData).toHaveBeenCalledWith(
			'/api/sources/7/events',
			20,
			null,
			expect.objectContaining({
				sort_direction: 'DESC'
			})
		);
	});

	it('loads another page when load more is clicked', async () => {
		fetchPaginatedData.mockReset();
		fetchPaginatedData
			.mockResolvedValueOnce({
				data: [
					{
						id: 11,
						source_id: 7,
						dedup_hash: 'abc123',
						method: 'POST',
						ingress_path: '/hook',
						query_params: {},
						raw_headers: {},
						body: btoa('payload'),
						body_content_type: 'text/plain',
						received_at: '2026-04-01T00:00:00.000Z'
					}
				],
				next_cursor: 'next-cursor',
				has_next: true
			})
			.mockResolvedValueOnce({
				data: [
					{
						id: 12,
						source_id: 7,
						dedup_hash: null,
						method: 'GET',
						ingress_path: '/hook-2',
						query_params: {},
						raw_headers: {},
						body: btoa('payload-2'),
						body_content_type: 'text/plain',
						received_at: '2026-04-02T00:00:00.000Z'
					}
				],
				next_cursor: null,
				has_next: false
			});

		render(ListEvents, {
			props: {
				sourceID: '7'
			}
		});

		await screen.findByRole('link', { name: 'Event ID: 11' });
		await fireEvent.click(screen.getByRole('button', { name: 'Load More' }));

		expect(await screen.findByRole('link', { name: 'Event ID: 12' })).toHaveAttribute(
			'href',
			'/sources/7/12'
		);
		expect(fetchPaginatedData).toHaveBeenNthCalledWith(
			2,
			'/api/sources/7/events',
			20,
			'next-cursor',
			expect.objectContaining({
				sort_direction: 'DESC'
			})
		);
	});
});
