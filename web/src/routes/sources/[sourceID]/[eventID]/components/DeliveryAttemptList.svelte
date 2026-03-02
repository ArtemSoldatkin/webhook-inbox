<script lang="ts">
	import { fetchPaginatedData } from '$lib/api';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import { parseDeliveryAttemptDTO } from '$lib/dtoParsers';
	import { type DeliveryAttemptDTO } from '$lib/types';

	export let sourceID: string;
	export let eventID: string;

	let data: DeliveryAttemptDTO[] = [];
	let loading = false;
	let error: string | null = null;

	let pageSize: number = 20;
	let nextCursor: string | null = null;
	let hasNext: boolean = false;

	async function fetchDeliveryAttempts() {
		loading = true;
		error = null;
		try {
			const result = await fetchPaginatedData(
				`/api/sources/${sourceID}/events/${eventID}/delivery-attempts`,
				pageSize,
				nextCursor
			);
			data = [...data, ...result.data.map(parseDeliveryAttemptDTO)];
			nextCursor = result.nextCursor;
			hasNext = result.hasNext;
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

	$: if (sourceID && eventID) {
		resetAndFetchDeliveryAttempts();
	}
</script>

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
			{#each data as attempt}
				<li>
					<section>
						<h3>Attempt ID: {attempt.ID}</h3>
						<p>Event ID: {attempt.EventID}</p>
						<p>Attempt Number: {attempt.AttemptNumber}</p>
						<p>Delivery State: {attempt.State}</p>
						<p>Status code: {attempt.StatusCode}</p>
						<p>Error Type: {attempt.ErrorType}</p>
						<p>Error Message: {attempt.ErrorMessage}</p>
						<p>
							Started at: {attempt.StartedAt ? new Date(attempt.StartedAt).toLocaleString() : 'N/A'}
						</p>
						<p>
							Finished at: {attempt.FinishedAt
								? new Date(attempt.FinishedAt).toLocaleString()
								: 'N/A'}
						</p>
						<p>Created at: {new Date(attempt.CreatedAt).toLocaleString()}</p>
						<p>
							Next attempt at:{' '}
							{attempt.NextAttemptAt ? new Date(attempt.NextAttemptAt).toLocaleString() : 'N/A'}
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
