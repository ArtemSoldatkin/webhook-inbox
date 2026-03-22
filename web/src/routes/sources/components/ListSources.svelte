<script lang="ts">
	import { resolve } from '$app/paths';
	import { fetchPaginatedData } from '$lib/api';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import { parseSourceDTO } from '$lib/dto-parsers';
	import type { SourceDTO } from '$lib/types';
	import { untrack } from 'svelte';

	/** Filters applied to the sources list. */
	type SourceFilters = {
		/** Requested page size for source results. */
		pageSize: number;
		/** Free-text query applied to the sources endpoint. */
		searchQuery: string;
		/** Selected source status filter. Supported values are 'active', 'paused', 'quarantined', 'disabled', and '*' for no filtering. */
		filterStatus: string;
		/** Sort order used for source results. */
		sortDirection: 'ASC' | 'DESC';
	};

	/** Loaded source rows for the current filters. */
	let data = $state<SourceDTO[]>([]);

	/** Tracks whether a source request is in flight. */
	let loading = $state(false);

	/** Holds the latest source loading error. */
	let error = $state<string | null>(null);

	/** Requested page size for the sources list. */
	let pageSize = $state(20);

	/** Cursor used to fetch the next page of sources. */
	let nextCursor = $state<string | null>(null);

	/** Indicates whether more source pages are available. */
	let hasNext = $state(false);

	/** Free-text query applied to the sources endpoint. */
	let searchQuery = $state('');

	/** Selected source status filter. */
	let filterStatus = $state('*');

	/** Available source status filter values. */
	const filterStatusOptions = ['active', 'paused', 'quarantined', 'disabled'];

	/** Sort order used for source results. */
	let sortDirection = $state<'ASC' | 'DESC'>('DESC');

	/**
	 * Collects the current filters into a single object for easier passing to fetch functions.
	 *
	 * @returns Current source filters.
	 */
	function getCurrentFilters(): SourceFilters {
		return {
			pageSize,
			searchQuery,
			filterStatus,
			sortDirection
		};
	}

	/** Loads the next page of sources for the active filters. */
	async function fetchSources(filters: SourceFilters): Promise<void> {
		loading = true;
		error = null;
		try {
			const urlSearchParams: Record<string, string> = {};
			if (filters.searchQuery) urlSearchParams.search = filters.searchQuery;
			if (filters.filterStatus) urlSearchParams.filter_status = filters.filterStatus;
			if (filters.sortDirection) urlSearchParams.sort_direction = filters.sortDirection;

			const result = await fetchPaginatedData(
				'/api/sources',
				filters.pageSize,
				nextCursor,
				urlSearchParams
			);
			data = [...data, ...result.data.map(parseSourceDTO)];
			nextCursor = result.next_cursor;
			hasNext = result.has_next;
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching sources:', err);
		} finally {
			loading = false;
		}
	}

	/** Clears the current source list and fetches the first page again. */
	async function resetAndFetchSources(filters: SourceFilters): Promise<void> {
		data = [];
		nextCursor = null;
		hasNext = false;
		await fetchSources(filters);
	}

	/** Refreshes the source list with the current filters when the user clicks the "Refresh" button. */
	function handleRefresh(): void {
		void resetAndFetchSources(getCurrentFilters());
	}

	/** Applies the current filters when the user clicks the "Search" button in the filter bar or presses Enter in the search input. */
	function handleSearch(): void {
		void resetAndFetchSources(getCurrentFilters());
	}

	/** Loads the next page of sources when the user clicks the "Load More" button. */
	function handleLoadMore(): void {
		void fetchSources(getCurrentFilters());
	}

	$effect(() => {
		const filters = getCurrentFilters();

		untrack(() => {
			void resetAndFetchSources(filters);
		});
	});
</script>

<FilterBar
	bind:searchQuery
	bind:filter={filterStatus}
	bind:sortDirection
	filterName="status"
	filterOptions={filterStatusOptions}
	onSearch={handleSearch}
/>
<button onclick={handleRefresh} disabled={loading}>Refresh Sources</button>
{#if loading}
	<p>Loading sources...</p>
{:else if error}
	<p class="error">Error: {error}</p>
{:else if data}
	{#if data.length > 0}
		<ul>
			{#each data as source (source.id)}
				<li>
					<section>
						<h2><a href={resolve(`/sources/${source.id}`)}>{source.id}</a></h2>
						<p>{source.description}</p>
						<p>{source.ingress_url}</p>
						<p>{source.egress_url}</p>
						<p>Static headers:</p>
						{#each Object.entries(source.static_headers ?? {}) as [key, value] (key)}
							<p>{key}: {value}</p>
						{/each}
						<p>{source.status}</p>
						<p>{source.status_reason}</p>
						<p>Created at: {new Date(source.created_at).toLocaleString()}</p>
						<p>Updated at: {new Date(source.updated_at).toLocaleString()}</p>
						<p>
							Disabled at: {source.disable_at
								? new Date(source.disable_at).toLocaleString()
								: 'N/A'}
						</p>
					</section>
				</li>
			{/each}
		</ul>
		{#if hasNext}
			<button onclick={handleLoadMore} disabled={loading}>Load More</button>
		{/if}
	{:else}
		<p>No sources found.</p>
	{/if}
{:else}
	<p>No sources found.</p>
{/if}
<PageSizeSelector bind:pageSize />
