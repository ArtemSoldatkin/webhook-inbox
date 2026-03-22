<script lang="ts">
	import { resolve } from '$app/paths';
	import { fetchPaginatedData } from '$lib/api';
	import DisplayMapOfStringArrays from '$lib/components/DisplayMapOfStringArrays.svelte';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import { parseEventDTO } from '$lib/dto-parsers';
	import type { EventDTO } from '$lib/types';
	import { untrack } from 'svelte';
	import BodyView from './BodyView.svelte';

	type Props = {
		/** Source id whose events are being loaded. */
		sourceID: string;
	};

	/** Source id whose events are being listed. */
	let { sourceID }: Props = $props();

	/** Loaded events for the current source. */
	let data = $state<EventDTO[]>([]);

	/** Tracks whether events are loading. */
	let loading = $state(false);

	/** Holds the latest event loading error. */
	let error = $state<string | null>(null);

	/** Requested page size for events. */
	let pageSize = $state(20);

	/** Cursor used to fetch the next page of events. */
	let nextCursor = $state<string | null>(null);

	/** Indicates whether more events are available. */
	let hasNext = $state(false);

	/** Free-text query applied to events. */
	let searchQuery = $state('');

	/** Sort order used for event results. */
	let sortDirection = $state<'ASC' | 'DESC'>('DESC');

	/**
	 * Builds the query parameters for the current event filters.
	 *
	 * @returns URL search params for the list request.
	 */
	function collectUrlSearchParams() {
		const params: Record<string, string> = {};
		if (searchQuery) {
			params.search = searchQuery;
		}
		if (sortDirection) {
			params.sort_direction = sortDirection;
		}
		return params;
	}

	/** Loads the next page of events for the active source. */
	async function fetchEvents() {
		loading = true;
		error = null;
		try {
			const urlSearchParams = collectUrlSearchParams();
			const result = await fetchPaginatedData(
				`/api/sources/${sourceID}/events`,
				pageSize,
				nextCursor,
				urlSearchParams
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

	/** Clears the current event list and reloads it from the first page. */
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

		untrack(() => {
			resetAndFetchEvents();
		});
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
