<script lang="ts">
	import DisplayMap from '$lib/components/DisplayMap.svelte';
	import DisplayMapOfStringArrays from '$lib/components/DisplayMapOfStringArrays.svelte';
	import type { EventDTO } from '$lib/types';

	export let sourceID: string;

	let data: EventDTO[] | null = null;
	let loading = false;
	let error: string | null = null;

	async function fetchEvents() {
		loading = true;
		error = null;
		try {
			const response = await fetch(`/api/sources/${sourceID}/events`);
			if (!response.ok) {
				throw new Error(`Failed to fetch events: ${response.statusText}`);
			}
			data = await response.json();
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching events:', err);
		} finally {
			loading = false;
		}
	}

	$: if (sourceID) {
		fetchEvents();
	}
</script>

<button on:click={fetchEvents}>Refresh Events</button>
{#if loading}
	<p>Loading events...</p>
{:else if error}
	<p class="error">Error: {error}</p>
{:else if data}
	{#if data.length === 0}
		<p>No events found for this source.</p>
	{:else}
		<ul>
			{#each data as event}
				<li>
					<section>
						<h3><a href={`/sources/${event.SourceID}/${event.ID}`}>Event ID: {event.ID}</a></h3>
						<p>Source ID: {event.SourceID}</p>
						<p>Method: {event.Method}</p>
						<DisplayMapOfStringArrays title="Query Parameters" data={event.QueryParams} />
						<DisplayMapOfStringArrays title="Raw Headers" data={event.RawHeaders} />
						<DisplayMap data={event.Body} />
					</section>
				</li>
			{/each}
		</ul>
	{/if}
{:else}
	<p>No events found for this source.</p>
{/if}
