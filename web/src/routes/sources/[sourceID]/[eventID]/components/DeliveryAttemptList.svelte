<script lang="ts">
	import { fetchPaginatedData } from '$lib/api';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { parseDeliveryAttemptDTO } from '$lib/dto-parsers';
	import { type DeliveryAttemptDTO } from '$lib/types';
	import { untrack } from 'svelte';

	type Props = {
		/** Source id that owns the current event. */
		sourceID: string;

		/** Event id whose delivery attempts are shown. */
		eventID: string;
	};

	/** Filters applied to the delivery attempts list. */
	type DeliveryAttemptFilters = {
		/** Source id that owns the current event. */
		sourceID: string;

		/** Event id whose delivery attempts are shown. */
		eventID: string;

		/** Requested page size for delivery attempt results. */
		pageSize: number;

		/** Free-text query applied to the delivery attempts endpoint. */
		searchQuery: string;

		/** Selected delivery state filter. Supported values are 'pending', 'in_flight', 'succeeded', 'failed', 'aborted', and '*' for no filtering. */
		filterState: string;

		/** Sort order used for delivery attempt results. */
		sortDirection: 'ASC' | 'DESC';
	};

	let { sourceID, eventID }: Props = $props();

	/** Loaded delivery attempts for the current event. */
	let data = $state<DeliveryAttemptDTO[]>([]);

	/** Tracks whether delivery attempts are loading. */
	let loading = $state(false);

	/** Holds the latest delivery attempt loading error. */
	let error = $state<string | null>(null);

	/** Requested page size for delivery attempts. */
	let pageSize = $state(20);

	/** Cursor used to fetch the next page of attempts. */
	let nextCursor = $state<string | null>(null);

	/** Indicates whether more delivery attempts are available. */
	let hasNext = $state(false);

	/** Free-text query applied to delivery attempts. */
	let searchQuery = $state('');

	/** Selected delivery state filter. */
	let filterState = $state('*');

	/** Available delivery attempt states for filtering. */
	const filterStateOptions = ['pending', 'in_flight', 'succeeded', 'failed', 'aborted'];

	/** Sort order used for delivery attempts. */
	let sortDirection = $state<'ASC' | 'DESC'>('DESC');

	/**
	 * Collects the current filters into a single object for easier passing to fetch functions.
	 *
	 * @returns Current delivery attempt filters.
	 */
	function getCurrentFilters(): DeliveryAttemptFilters {
		return {
			sourceID,
			eventID,
			pageSize,
			searchQuery,
			filterState,
			sortDirection
		};
	}

	/** Loads the next page of delivery attempts.
	 *
	 * @param filters - Filters to apply when fetching delivery attempts.
	 */
	async function fetchDeliveryAttempts(filters: DeliveryAttemptFilters): Promise<void> {
		loading = true;
		error = null;
		try {
			const urlSearchParams: Record<string, string> = {};
			if (filters.searchQuery) urlSearchParams.search = filters.searchQuery;
			if (filters.filterState) urlSearchParams.filter_state = filters.filterState;
			if (filters.sortDirection) urlSearchParams.sort_direction = filters.sortDirection;

			const result = await fetchPaginatedData(
				`/api/sources/${filters.sourceID}/events/${filters.eventID}/delivery-attempts`,
				filters.pageSize,
				nextCursor,
				urlSearchParams
			);
			data = [...data, ...result.data.map(parseDeliveryAttemptDTO)];
			nextCursor = result.next_cursor;
			hasNext = result.has_next;
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching event details:', err);
		} finally {
			loading = false;
		}
	}

	/** Resets delivery attempt pagination and fetches from the start.
	 *
	 * @param filters - Filters to apply when fetching delivery attempts.
	 */
	async function resetAndFetchDeliveryAttempts(filters: DeliveryAttemptFilters): Promise<void> {
		data = [];
		nextCursor = null;
		hasNext = false;
		await fetchDeliveryAttempts(filters);
	}

	/** Refreshes the delivery attempt list with the current filters when the user clicks the "Refresh" button. */
	function handleRefresh(): void {
		void resetAndFetchDeliveryAttempts(getCurrentFilters());
	}

	/** Loads the next page of delivery attempts when the user clicks the "Load More" button. */
	function handleLoadMore(): void {
		void fetchDeliveryAttempts(getCurrentFilters());
	}

	$effect(() => {
		const filters = getCurrentFilters();

		untrack(() => {
			void resetAndFetchDeliveryAttempts(filters);
		});
	});
</script>

<FilterBar
	bind:searchQuery
	bind:filter={filterState}
	bind:sortDirection
	filterName="state"
	filterOptions={filterStateOptions}
/>
<Button onclick={handleRefresh} disabled={loading}>Refresh Delivery Attempts</Button>
<h3>Delivery Attempts</h3>
{#if loading}
	<p>Loading event details...</p>
{:else if error}
	<p class="error">Error: {error}</p>
{:else if data}
	{#if data.length === 0}
		<p>No delivery attempts found for this event.</p>
	{:else}
		<ul>
			{#each data as attempt (attempt.id)}
				<li>
					<section>
						<h3>Attempt ID: {attempt.id}</h3>
						<p>Event ID: {attempt.event_id}</p>
						<p>Attempt Number: {attempt.attempt_number}</p>
						<p>Delivery State: {attempt.state}</p>
						<p>Status code: {attempt.status_code}</p>
						<p>Error Type: {attempt.error_type}</p>
						<p>Error Message: {attempt.error_message}</p>
						<p>
							Started at: {attempt.started_at
								? new Date(attempt.started_at).toLocaleString()
								: 'N/A'}
						</p>
						<p>
							Finished at: {attempt.finished_at
								? new Date(attempt.finished_at).toLocaleString()
								: 'N/A'}
						</p>
						<p>Created at: {new Date(attempt.created_at).toLocaleString()}</p>
						<p>
							Next attempt at:
							{attempt.next_attempt_at ? new Date(attempt.next_attempt_at).toLocaleString() : 'N/A'}
						</p>
					</section>
				</li>
			{/each}
		</ul>
		{#if hasNext}
			<Button onclick={handleLoadMore} disabled={loading}>Load More</Button>
		{/if}
	{/if}
{:else}
	<p>No details found for this event.</p>
{/if}
<PageSizeSelector bind:pageSize />
