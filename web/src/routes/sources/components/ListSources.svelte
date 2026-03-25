<script lang="ts">
	import { resolve } from '$app/paths';
	import { fetchPaginatedData } from '$lib/api';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import KayValueList from '$lib/components/ui/KayValueList.svelte';
	import Link from '$lib/components/ui/Link.svelte';
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

<section class="rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8">
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
			<div>
				<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">Sources</p>
				<h2 class="mt-4 text-3xl font-semibold tracking-tight text-fg">Manage registered endpoints</h2>
				<p class="mt-3 max-w-2xl text-sm leading-6 text-muted sm:text-base">
					Browse created sources, inspect forwarding settings, and drill into event history for each
					endpoint.
				</p>
			</div>
			<Button onclick={handleRefresh} disabled={loading} variant="secondary">Refresh Sources</Button>
		</div>

		<div class="rounded-lg border border-border-muted bg-elevated p-4">
			<FilterBar
				bind:searchQuery
				bind:filter={filterStatus}
				bind:sortDirection
				filterName="status"
				filterOptions={filterStatusOptions}
			/>
		</div>

		{#if loading}
			<div class="rounded-md border border-border-muted bg-elevated px-4 py-6 text-sm text-muted">
				Loading sources...
			</div>
		{:else if error}
			<div class="rounded-md border border-error bg-surface px-4 py-3 text-sm text-error">
				Error: {error}
			</div>
		{:else if data.length > 0}
			<ul class="grid gap-4">
				{#each data as source (source.id)}
					<li>
						<article class="rounded-lg border border-border bg-elevated p-5 shadow-sm">
							<div class="flex flex-col gap-5">
								<div class="flex flex-col gap-2">
									<div class="flex flex-wrap items-center gap-3">
										<h3 class="text-xl font-semibold tracking-tight text-fg">
											<Link href={resolve(`/sources/${source.id}`)} variant="inline">{source.id}</Link>
										</h3>
										<span
											class="inline-flex w-fit rounded-full border border-border bg-surface px-3 py-1 text-xs font-medium uppercase tracking-[0.12em] text-muted"
										>
											{source.status}
										</span>
									</div>
									<p class="text-sm leading-6 text-muted">
										{source.description || 'No description provided.'}
									</p>
								</div>

								<div class="grid gap-4 sm:grid-cols-2">
									<div class="rounded-md border border-border-muted bg-surface p-4">
										<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
											Ingress URL
										</p>
										<p class="mt-2 break-all text-sm leading-6 text-fg">{source.ingress_url}</p>
									</div>
									<div class="rounded-md border border-border-muted bg-surface p-4">
										<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
											Egress URL
										</p>
										<p class="mt-2 break-all text-sm leading-6 text-fg">{source.egress_url}</p>
									</div>
								</div>

								<div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,0.9fr)]">
									<div class="rounded-md border border-border-muted bg-surface p-4">
										<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
											Static headers
										</p>
										{#if Object.keys(source.static_headers ?? {}).length > 0}
											<div class="mt-3 flex flex-col gap-2">
												{#each Object.entries(source.static_headers ?? {}) as [key, value] (key)}
													<div
														class="flex flex-col gap-1 rounded-md border border-border-muted bg-elevated px-3 py-2 sm:flex-row sm:items-start sm:justify-between sm:gap-4"
													>
														<span class="text-sm font-medium text-fg">{key}</span>
														<span class="break-all text-sm text-muted">{value}</span>
													</div>
												{/each}
											</div>
										{:else}
											<p class="mt-2 text-sm text-muted">No static headers configured.</p>
										{/if}
									</div>

									<div class="rounded-md border border-border-muted bg-surface p-4">
										<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
											Metadata
										</p>
										<KayValueList
											items={[
												{ label: 'Status reason', value: source.status_reason },
												{ label: 'Created at', value: new Date(source.created_at).toLocaleString() },
												{ label: 'Updated at', value: new Date(source.updated_at).toLocaleString() },
												{
													label: 'Disabled at',
													value: source.disable_at
														? new Date(source.disable_at).toLocaleString()
														: 'N/A'
												}
											]}
										/>
									</div>
								</div>
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
		{:else}
			<div class="rounded-md border border-border-muted bg-elevated px-4 py-6 text-sm text-muted">
				No sources found.
			</div>
		{/if}

		<div class="border-t border-border-muted pt-4">
			<PageSizeSelector bind:pageSize />
		</div>
	</div>
</section>
