<script lang="ts">
	import { parseDeliveryAttemptDTO } from '$lib/dtoParsers';
	import { type DeliveryAttemptDTO } from '$lib/types';

	export let sourceID: string;
	export let eventID: string;

	let data: DeliveryAttemptDTO[] | null = null;
	let loading = false;
	let error: string | null = null;

	async function fetchDeliveryAttempts() {
		loading = true;
		error = null;
		try {
			const response = await fetch(`/api/sources/${sourceID}/events/${eventID}/delivery_attempts`);
			if (!response.ok) {
				throw new Error(`Failed to fetch event details: ${response.statusText}`);
			}
			const rawData = await response.json();
			data = rawData.map(parseDeliveryAttemptDTO);
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching event details:', err);
		} finally {
			loading = false;
		}
	}

	$: if (sourceID && eventID) {
		fetchDeliveryAttempts();
	}
</script>

<h3>Delivery Attempts</h3>
{#if loading}
	<p>Loading event details...</p>
{:else if error}
	<p class="error">Error: {error}</p>
{:else if data}
	{#if data.length === 0}
		<p>No delivery attempts found for this event.</p>
	{:else}
		{#each data as attempt}
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
					Finished at: {attempt.FinishedAt ? new Date(attempt.FinishedAt).toLocaleString() : 'N/A'}
				</p>
				<p>Created at: {new Date(attempt.CreatedAt).toLocaleString()}</p>
				<p>
					Next attempt at:{' '}
					{attempt.NextAttemptAt ? new Date(attempt.NextAttemptAt).toLocaleString() : 'N/A'}
				</p>
			</section>
		{/each}
	{/if}
{:else}
	<p>No details found for this event.</p>
{/if}
