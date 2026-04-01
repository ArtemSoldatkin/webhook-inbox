import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import DeliveryAttemptList from './DeliveryAttemptList.svelte';

const { fetchPaginatedData } = vi.hoisted(() => ({
	fetchPaginatedData: vi.fn()
}));

vi.mock('$lib/api', () => ({
	fetchPaginatedData
}));

describe('DeliveryAttemptList', () => {
	it('loads and renders delivery attempts on mount', async () => {
		fetchPaginatedData.mockReset();
		fetchPaginatedData.mockResolvedValueOnce({
			data: [
				{
					id: 101,
					event_id: 44,
					attempt_number: 1,
					state: 'failed',
					status_code: 503,
					error_type: 'network_error',
					error_message: 'connection timed out',
					started_at: '2026-04-01T00:00:00.000Z',
					finished_at: '2026-04-01T00:00:05.000Z',
					created_at: '2026-04-01T00:00:00.000Z',
					next_attempt_at: '2026-04-01T00:05:00.000Z'
				}
			],
			next_cursor: null,
			has_next: false
		});

		render(DeliveryAttemptList, {
			props: {
				sourceID: '12',
				eventID: '44'
			}
		});

		expect(await screen.findByRole('heading', { name: 'Attempt ID: 101' })).toBeInTheDocument();
		expect(
			screen.getByRole('heading', { name: 'Attempt ID: 101' }).closest('article')
		).toHaveTextContent('failed');
		expect(screen.getByText('connection timed out')).toBeInTheDocument();
		expect(fetchPaginatedData).toHaveBeenCalledWith(
			'/api/sources/12/events/44/delivery-attempts',
			20,
			null,
			expect.objectContaining({
				filter_state: '*',
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
						id: 101,
						event_id: 44,
						attempt_number: 1,
						state: 'failed',
						status_code: 503,
						error_type: 'network_error',
						error_message: 'connection timed out',
						started_at: '2026-04-01T00:00:00.000Z',
						finished_at: '2026-04-01T00:00:05.000Z',
						created_at: '2026-04-01T00:00:00.000Z',
						next_attempt_at: '2026-04-01T00:05:00.000Z'
					}
				],
				next_cursor: 'next-cursor',
				has_next: true
			})
			.mockResolvedValueOnce({
				data: [
					{
						id: 102,
						event_id: 44,
						attempt_number: 2,
						state: 'succeeded',
						status_code: 200,
						error_type: null,
						error_message: null,
						started_at: '2026-04-01T00:05:00.000Z',
						finished_at: '2026-04-01T00:05:02.000Z',
						created_at: '2026-04-01T00:05:00.000Z',
						next_attempt_at: null
					}
				],
				next_cursor: null,
				has_next: false
			});

		render(DeliveryAttemptList, {
			props: {
				sourceID: '12',
				eventID: '44'
			}
		});

		await screen.findByRole('heading', { name: 'Attempt ID: 101' });
		await fireEvent.click(screen.getByRole('button', { name: 'Load More' }));

		expect(await screen.findByRole('heading', { name: 'Attempt ID: 102' })).toBeInTheDocument();
		expect(fetchPaginatedData).toHaveBeenNthCalledWith(
			2,
			'/api/sources/12/events/44/delivery-attempts',
			20,
			'next-cursor',
			expect.objectContaining({
				filter_state: '*',
				sort_direction: 'DESC'
			})
		);
	});

	it('shows an error alert when fetching attempts fails', async () => {
		fetchPaginatedData.mockReset();
		fetchPaginatedData.mockRejectedValueOnce(new Error('Request failed'));

		render(DeliveryAttemptList, {
			props: {
				sourceID: '12',
				eventID: '44'
			}
		});

		expect(await screen.findByText('Request failed')).toBeInTheDocument();
		expect(
			screen.queryByText('No delivery attempts found for this event.')
		).not.toBeInTheDocument();
	});
});
