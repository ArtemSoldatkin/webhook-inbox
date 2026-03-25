<script lang="ts">
	import { resolve } from '$app/paths';
	import { fetchPaginatedData } from '$lib/api';
	import DisplayMapOfStringArrays from '$lib/components/DisplayMapOfStringArrays.svelte';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Link from '$lib/components/ui/Link.svelte';
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

<section class="rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8">
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
			<div>
				<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">Events</p>
				<h3 class="mt-4 text-2xl font-semibold tracking-tight text-fg">Captured traffic for this source</h3>
				<p class="mt-3 max-w-2xl text-sm leading-6 text-muted sm:text-base">
					Inspect recorded requests, query parameters, headers, and request bodies in arrival order.
				</p>
			</div>
			<Button onclick={handleRefresh} disabled={loading} variant="secondary">Refresh Events</Button>
		</div>

		<div class="rounded-lg border border-border-muted bg-elevated p-4">
			<FilterBar bind:searchQuery bind:sortDirection />
		</div>

		{#if loading}
			<div class="rounded-md border border-border-muted bg-elevated px-4 py-6 text-sm text-muted">
				Loading events...
			</div>
		{:else if error}
			<div class="rounded-md border border-error bg-surface px-4 py-3 text-sm text-error">
				Error: {error}
			</div>
		{:else if data.length === 0}
			<div class="rounded-md border border-border-muted bg-elevated px-4 py-6 text-sm text-muted">
				No events found for this source.
			</div>
		{:else}
			<ul class="grid gap-4">
				{#each data as event (event.id)}
					<li>
						<article class="rounded-lg border border-border bg-elevated p-5 shadow-sm">
							<div class="flex flex-col gap-5">
								<div class="flex flex-col gap-2">
									<div class="flex flex-wrap items-center gap-3">
										<h4 class="text-xl font-semibold tracking-tight text-fg">
											<Link href={resolve(`/sources/${event.source_id}/${event.id}`)} variant="inline">
												Event ID: {event.id}
											</Link>
										</h4>
										<span
											class="inline-flex w-fit rounded-full border border-border bg-surface px-3 py-1 text-xs font-medium uppercase tracking-[0.12em] text-muted"
										>
											{event.method}
										</span>
									</div>
									<div class="grid gap-3 text-sm sm:grid-cols-2">
										<div class="rounded-md border border-border-muted bg-surface px-3 py-2">
											<span class="text-muted">Source ID</span>
											<p class="mt-1 break-all text-fg">{event.source_id}</p>
										</div>
										<div class="rounded-md border border-border-muted bg-surface px-3 py-2">
											<span class="text-muted">Deduplication hash</span>
											<p class="mt-1 break-all text-fg">{event.dedup_hash ?? 'N/A'}</p>
										</div>
									</div>
								</div>

								<div class="grid gap-4 lg:grid-cols-2">
									<DisplayMapOfStringArrays title="Query Parameters" data={event.query_params ?? {}} />
									<DisplayMapOfStringArrays title="Raw Headers" data={event.raw_headers ?? {}} />
								</div>

								<BodyView body={event.body} contentType={event.body_content_type} />
							</div>
						</article>
					</li>
				{/each}
			</ul>

			{#if hasNext}
				<div class="flex justify-center pt-2">
					<Button onclick={handleLoadMore} disabled={loading} variant="secondary">Load More</Button>
				</div>
			{/if}
		{/if}

		<div class="border-t border-border-muted pt-4">
			<PageSizeSelector bind:pageSize />
		</div>
	</div>
</section>
