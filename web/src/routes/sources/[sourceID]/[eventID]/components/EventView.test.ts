import { render, screen } from '@testing-library/svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';
import EventView from './EventView.svelte';

vi.mock('$lib/api', () => ({
	fetchPaginatedData: vi.fn().mockResolvedValue({
		data: [],
		next_cursor: null,
		has_next: false
	})
}));

describe('EventView', () => {
	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('loads and renders event details', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({
					id: 44,
					source_id: 12,
					dedup_hash: 'abc123',
					method: 'POST',
					ingress_path: '/hook',
					query_params: { tenant: ['acme'] },
					raw_headers: { 'content-type': ['application/json'] },
					body: btoa('plain payload'),
					body_content_type: 'text/plain',
					received_at: '2026-04-01T00:00:00.000Z'
				})
			})
		);

		render(EventView, {
			props: {
				sourceID: '12',
				eventID: '44'
			}
		});

		expect(await screen.findByRole('heading', { name: 'Event ID: 44' })).toBeInTheDocument();
		expect(screen.getByText('abc123')).toBeInTheDocument();
		expect(screen.getByText('tenant')).toBeInTheDocument();
		expect(
			screen.getByRole('heading', { name: 'Delivery history for this event' })
		).toBeInTheDocument();
	});
});
