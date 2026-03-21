<script lang="ts">
	import { resolve } from '$app/paths';
	import { fetchPaginatedData } from '$lib/api';
	import DisplayMapOfStringArrays from '$lib/components/DisplayMapOfStringArrays.svelte';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import { parseEventDTO } from '$lib/dtoParsers';
	import type { EventDTO } from '$lib/types';
	import BodyView from './BodyView.svelte';

	type Props = {
		sourceID: string;
	};

	let { sourceID }: Props = $props();

	let data = $state<EventDTO[]>([]);
	let loading = $state(false);
	let error = $state<string | null>(null);

	let pageSize = $state(20);
	let nextCursor = $state<string | null>(null);
	let hasNext = $state(false);

	let searchQuery = $state('');
	let sortDirection = $state<'ASC' | 'DESC'>('DESC');

	async function fetchEvents() {
		loading = true;
		error = null;
		try {
			const result = await fetchPaginatedData<EventDTO>(
				`/api/sources/${sourceID}/events`,
				pageSize,
				nextCursor,
				{
					search: searchQuery,
					sort_direction: sortDirection
				}
			);
			data = [...data, ...result.data.map(parseEventDTO)];
			nextCursor = result.next_cursor;
			hasNext = result.has_next;
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching events:', err);
		} finally {
			loading = false;
		}
	}

	async function resetAndFetchEvents() {
		data = [];
		nextCursor = null;
		hasNext = false;
		await fetchEvents();
	}

	$effect(() => {
		sourceID;
		pageSize;
		sortDirection;
		resetAndFetchEvents();
	});
</script>

<FilterBar bind:searchQuery bind:sortDirection onSearch={resetAndFetchEvents} />
<button onclick={resetAndFetchEvents} disabled={loading}>Refresh Events</button>
{#if loading}
	<p>Loading events...</p>
{:else if error}
	<p class="error">Error: {error}</p>
{:else if data}
	{#if data.length === 0}
		<p>No events found for this source.</p>
	{:else}
		<ul>
			{#each data as event (event.id)}
				<li>
					<section>
						<h3>
							<a href={resolve(`/sources/${event.source_id}/${event.id}`)}>Event ID: {event.id}</a>
						</h3>
						<p>Source ID: {event.source_id}</p>
						<p>Deduplication Hash: {event.dedup_hash ?? 'N/A'}</p>
						<p>Method: {event.method}</p>
						<DisplayMapOfStringArrays title="Query Parameters" data={event.query_params ?? {}} />
						<DisplayMapOfStringArrays title="Raw Headers" data={event.raw_headers ?? {}} />
						<BodyView body={event.body} contentType={event.body_content_type} />
					</section>
				</li>
			{/each}
		</ul>
		{#if hasNext}
			<button onclick={fetchEvents} disabled={loading}>Load More</button>
		{/if}
	{/if}
{:else}
	<p>No events found for this source.</p>
{/if}
<PageSizeSelector bind:pageSize />
