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

	/** Filters applied to the events list. */
	type EventFilters = {
		/** Source id whose events are being loaded. */
		sourceID: string;
		/** Requested page size for event results. */
		pageSize: number;
		/** Free-text query applied to the events endpoint. */
		searchQuery: string;
		/** Sort order used for event results. */
		sortDirection: 'ASC' | 'DESC';
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
	 * Collects the current filters into a single object for easier passing to fetch functions.
	 *
	 * @returns Current event filters.
	 */
	function getCurrentFilters(): EventFilters {
		return {
			sourceID,
			pageSize,
			searchQuery,
			sortDirection
		};
	}

	/**
	 * Loads the next page of events for the active source.
	 *
	 * @param filters - Filters to apply when fetching events.
	 */
	async function fetchEvents(filters: EventFilters): Promise<void> {
		loading = true;
		error = null;
		try {
			const urlSearchParams: Record<string, string> = {};
			if (filters.searchQuery) urlSearchParams.search = filters.searchQuery;
			if (filters.sortDirection) urlSearchParams.sort_direction = filters.sortDirection;

			const result = await fetchPaginatedData(
				`/api/sources/${filters.sourceID}/events`,
				filters.pageSize,
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

	/** Clears the current event list and reloads it from the first page.
	 *
	 * @param filters - Filters to apply when fetching events.
	 */
	async function resetAndFetchEvents(filters: EventFilters): Promise<void> {
		data = [];
		nextCursor = null;
		hasNext = false;
		await fetchEvents(filters);
	}

	/** Refreshes the event list with the current filters when the user clicks the "Refresh" button. */
	function handleRefresh(): void {
		void resetAndFetchEvents(getCurrentFilters());
	}

	/** Applies the current filters when the user triggers a search from the filter bar (for example, by clicking the "Search" button). */
	function handleSearch(): void {
		void resetAndFetchEvents(getCurrentFilters());
	}

	/** Loads the next page of events when the user clicks the "Load More" button. */
	function handleLoadMore(): void {
		void fetchEvents(getCurrentFilters());
	}

	$effect(() => {
		const filters = getCurrentFilters();

		untrack(() => {
			void resetAndFetchEvents(filters);
		});
	});
</script>

<FilterBar bind:searchQuery bind:sortDirection onSearch={handleSearch} />
<button onclick={handleRefresh} disabled={loading}>Refresh Events</button>
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
			<button onclick={handleLoadMore} disabled={loading}>Load More</button>
		{/if}
	{/if}
{:else}
	<p>No events found for this source.</p>
{/if}
<PageSizeSelector bind:pageSize />
