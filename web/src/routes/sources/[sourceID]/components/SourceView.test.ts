import { render, screen } from '@testing-library/svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';
import SourceView from './SourceView.svelte';

vi.mock('$app/paths', () => ({
	resolve: (path: string) => path
}));

vi.mock('$lib/api', () => ({
	fetchPaginatedData: vi.fn().mockResolvedValue({
		data: [],
		next_cursor: null,
		has_next: false
	})
}));

describe('SourceView', () => {
	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('loads and renders source details', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn().mockResolvedValue({
				ok: true,
				json: async () => ({
					id: 12,
					public_id: 'src_12',
					ingress_url: 'http://localhost:3000/api/v1/ingest/src_12',
					egress_url: 'https://example.com/out',
					static_headers: { Authorization: 'Bearer token' },
					status: 'active',
					status_reason: 'ready',
					description: 'Payments source',
					created_at: '2026-04-01T00:00:00.000Z',
					updated_at: '2026-04-02T00:00:00.000Z',
					disable_at: null
				})
			})
		);

		render(SourceView, {
			props: {
				sourceID: '12'
			}
		});

		expect(await screen.findByText('Payments source')).toBeInTheDocument();
		expect(screen.getByRole('heading', { name: 'Send a test webhook' })).toBeInTheDocument();
		expect(screen.getByRole('heading', { name: 'Captured traffic for this source' })).toBeInTheDocument();
	});
});
