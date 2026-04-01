import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';
import SourceCard from './SourceCard.svelte';

vi.mock('$app/paths', () => ({
	resolve: (path: string) => path
}));

function createSource(overrides: Record<string, unknown> = {}) {
	return {
		id: 7,
		public_id: 'src_7',
		ingress_url: 'http://localhost:3000/api/v1/ingest/src_7',
		egress_url: 'https://example.com/out',
		static_headers: { Authorization: 'Bearer token' },
		status: 'active',
		status_reason: 'ready',
		description: 'Primary source',
		created_at: new Date('2026-04-01T00:00:00.000Z'),
		updated_at: new Date('2026-04-02T00:00:00.000Z'),
		disable_at: null,
		...overrides
	};
}

describe('SourceCard', () => {
	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('renders the source details and optional link state', () => {
		render(SourceCard, {
			props: {
				source: createSource(),
				idAsLink: true
			}
		});

		expect(screen.getByRole('link', { name: '7' })).toHaveAttribute('href', '/sources/7');
		expect(screen.getByText('Primary source')).toBeInTheDocument();
		expect(screen.getByText('http://localhost:3000/api/v1/ingest/src_7')).toBeInTheDocument();
		expect(screen.getByText('https://example.com/out')).toBeInTheDocument();
		expect(screen.getByText('Authorization')).toBeInTheDocument();
		expect(screen.getByText('Bearer token')).toBeInTheDocument();
		expect(screen.getByText('active')).toBeInTheDocument();
	});

	it('updates the source status and calls the refresh callback', async () => {
		const fetchMock = vi.fn().mockResolvedValue({
			ok: true,
			json: async () => ({})
		});
		const onStatusUpdate = vi.fn();
		vi.stubGlobal('fetch', fetchMock);

		render(SourceCard, {
			props: {
				source: createSource(),
				onStatusUpdate
			}
		});

		await fireEvent.click(screen.getByRole('button', { name: 'Edit source status' }));
		await fireEvent.change(screen.getByRole('combobox'), {
			target: { value: 'paused' }
		});
		await fireEvent.click(screen.getByRole('button', { name: 'Save source status' }));

		await waitFor(() => {
			expect(fetchMock).toHaveBeenCalledWith(
				'/api/sources/7/status',
				expect.objectContaining({
					method: 'PUT',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({
						status: 'paused',
						status_reason: 'Status changed to paused via UI'
					})
				})
			);
		});
		await waitFor(() => {
			expect(onStatusUpdate).toHaveBeenCalledTimes(1);
		});
		expect(screen.getByText('paused')).toBeInTheDocument();
	});

	it('shows an error when updating the source status fails', async () => {
		const fetchMock = vi.fn().mockResolvedValue({
			ok: false,
			json: async () => ({ message: 'Failed to update source status' })
		});
		vi.stubGlobal('fetch', fetchMock);
		const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});

		render(SourceCard, {
			props: {
				source: createSource({ static_headers: {} })
			}
		});

		expect(screen.getByText('No static headers configured.')).toBeInTheDocument();

		await fireEvent.click(screen.getByRole('button', { name: 'Edit source status' }));
		await fireEvent.change(screen.getByRole('combobox'), {
			target: { value: 'disabled' }
		});
		await fireEvent.click(screen.getByRole('button', { name: 'Save source status' }));

		expect(
			await screen.findByText('Failed to update source status. Please try again.')
		).toBeInTheDocument();
		expect(screen.getByText('No static headers configured.')).toBeInTheDocument();
		expect(screen.queryByText('disabled')).not.toBeInTheDocument();
		expect(consoleError).toHaveBeenCalled();
	});
});
