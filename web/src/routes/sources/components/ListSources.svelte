<script lang="ts">
	import { fetchPaginatedData } from '$lib/api';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { parseSourceDTO } from '$lib/dto-parsers';
	import type { SourceDTO } from '$lib/types';
	import { cx } from '$lib/utils/cx';
	import { untrack } from 'svelte';
	import SectionHeader from '$lib/components/ui/SectionHeader.svelte';
	import SourceCard from './SourceCard.svelte';

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

	type Props = {
		/** Additional CSS classes to apply to the root element of this component. */
		class?: string;
	};

	let { class: className }: Props = $props();

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

	/**
	 * Loads the next page of sources for the active filters.
	 *
	 * @param filters - Filters to apply when fetching sources.
	 */
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

	/** Clears the current source list and fetches the first page again.
	 *
	 * @param filters - Filters to apply when fetching sources.
	 */
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

<section class={cx('rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8', className)}>
	<div class="flex flex-col gap-6">
		<SectionHeader
			eyebrow="Sources"
			title="Manage registered endpoints"
			description="Browse created sources, inspect forwarding settings, and drill into event history for each endpoint."
		>
			{#snippet actions()}
				<Button onclick={handleRefresh} disabled={loading} variant="secondary">
					Refresh Sources
				</Button>
			{/snippet}
		</SectionHeader>

		<FilterBar
			bind:searchQuery
			bind:filter={filterStatus}
			bind:sortDirection
			filterName="status"
			filterOptions={filterStatusOptions}
		/>

		{#if loading}
			<Alert>Loading sources...</Alert>
		{:else if error}
			<Alert variant="error" title="Error" class="bg-surface">{error}</Alert>
		{:else if data.length > 0}
			<ul class="grid gap-4">
				{#each data as source (source.id)}
					<li>
						<SourceCard {source} idAsLink />
					</li>
				{/each}
			</ul>
			{#if hasNext}
				<div class="flex justify-center pt-2">
					<Button onclick={handleLoadMore} disabled={loading} variant="secondary">Load More</Button>
				</div>
			{/if}
		{:else}
			<Alert>No sources found.</Alert>
		{/if}

		<div class="border-t border-border-muted pt-4">
			<PageSizeSelector bind:pageSize />
		</div>
	</div>
</section>
