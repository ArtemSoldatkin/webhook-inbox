<script lang="ts">
	import { fetchPaginatedData } from '$lib/api';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import { parseDeliveryAttemptDTO } from '$lib/dtoParsers';
	import { type DeliveryAttemptDTO } from '$lib/types';

	type Props = {
		sourceID: string;
		eventID: string;
	};

	let { sourceID, eventID }: Props = $props();

	let data = $state<DeliveryAttemptDTO[]>([]);
	let loading = $state(false);
	let error = $state<string | null>(null);

	let pageSize = $state(20);
	let nextCursor = $state<string | null>(null);
	let hasNext = $state(false);

	let searchQuery = $state('');

	let filterState = $state('*');
	const filterStateOptions = ['pending', 'in_flight', 'succeeded', 'failed', 'aborted'];

	let sortDirection = $state<'ASC' | 'DESC'>('DESC');

	async function fetchDeliveryAttempts() {
		loading = true;
		error = null;
		try {
			const result = await fetchPaginatedData(
				`/api/sources/${sourceID}/events/${eventID}/delivery-attempts`,
				pageSize,
				nextCursor,
				{
					search: searchQuery,
					filter_state: filterState,
					sort_direction: sortDirection
				}
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
		resetAndFetchDeliveryAttempts();
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
<button on:click={resetAndFetchDeliveryAttempts} disabled={loading}>Refresh Events</button>
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
			<button on:click={fetchDeliveryAttempts} disabled={loading}>Load More</button>
		{/if}
	{/if}
{:else}
	<p>No details found for this event.</p>
{/if}
<PageSizeSelector bind:pageSize />
