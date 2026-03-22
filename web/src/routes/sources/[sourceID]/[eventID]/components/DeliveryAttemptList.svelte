<script lang="ts">
	import { fetchPaginatedData } from '$lib/api';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import { parseDeliveryAttemptDTO } from '$lib/dto-parsers';
	import { type DeliveryAttemptDTO } from '$lib/types';
	import { untrack } from 'svelte';

	type Props = {
		/** Source id that owns the current event. */
		sourceID: string;

		/** Event id whose delivery attempts are shown. */
		eventID: string;
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
	 * Builds the query parameters for the delivery attempt request.
	 *
	 * @returns URL search params for the API call.
	 */
	function collectUrlSearchParams() {
		const params: Record<string, string> = {};
		if (searchQuery) {
			params.search = searchQuery;
		}
		if (filterState) {
			params.filter_state = filterState;
		}
		if (sortDirection) {
			params.sort_direction = sortDirection;
		}
		return params;
	}

	/** Loads the next page of delivery attempts. */
	async function fetchDeliveryAttempts() {
		loading = true;
		error = null;
		try {
			const urlSearchParams = collectUrlSearchParams();
			const result = await fetchPaginatedData(
				`/api/sources/${sourceID}/events/${eventID}/delivery-attempts`,
				pageSize,
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

	/** Resets delivery attempt pagination and fetches from the start. */
	async function resetAndFetchDeliveryAttempts() {
		data = [];
		nextCursor = null;
		hasNext = false;
		await fetchDeliveryAttempts();
	}

	$effect(() => {
		sourceID;
		eventID;
		pageSize;
		filterState;
		sortDirection;

		untrack(() => {
			resetAndFetchDeliveryAttempts();
		});
	});
</script>

<FilterBar
	bind:searchQuery
	bind:filter={filterState}
	bind:sortDirection
	filterName="state"
	filterOptions={filterStateOptions}
	onSearch={resetAndFetchDeliveryAttempts}
/>
<button onclick={resetAndFetchDeliveryAttempts} disabled={loading}>Refresh Events</button>
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
			<button onclick={fetchDeliveryAttempts} disabled={loading}>Load More</button>
		{/if}
	{/if}
{:else}
	<p>No details found for this event.</p>
{/if}
<PageSizeSelector bind:pageSize />
