<script lang="ts">
	import { resolve } from '$app/paths';
	import { fetchPaginatedData } from '$lib/api';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import { parseSourceDTO } from '$lib/dto-parsers';
	import type { SourceDTO } from '$lib/types';
	import { untrack } from 'svelte';

	let data = $state<SourceDTO[]>([]);
	let loading = $state(false);
	let error = $state<string | null>(null);

	let pageSize = $state(20);
	let nextCursor = $state<string | null>(null);
	let hasNext = $state(false);

	let searchQuery = $state('');

	let filterStatus = $state('*');
	const filterStatusOptions = ['active', 'paused', 'quarantined', 'disabled'];

	let sortDirection = $state<'ASC' | 'DESC'>('DESC');

	function collectUrlSearchParams() {
		const params: Record<string, string> = {};
		if (searchQuery) {
			params.search = searchQuery;
		}
		if (filterStatus) {
			params.filter_status = filterStatus;
		}
		if (sortDirection) {
			params.sort_direction = sortDirection;
		}
		return params;
	}

	async function fetchSources() {
		loading = true;
		error = null;
		try {
			const urlSearchParams = collectUrlSearchParams();
			const result = await fetchPaginatedData(
				'/api/sources',
				pageSize,
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

	async function resetAndFetchSources() {
		data = [];
		nextCursor = null;
		hasNext = false;
		await fetchSources();
	}

	$effect(() => {
		pageSize;
		filterStatus;
		sortDirection;

		untrack(() => {
			resetAndFetchSources();
		});
	});
</script>

<FilterBar
	bind:searchQuery
	bind:filter={filterStatus}
	bind:sortDirection
	filterName="status"
	filterOptions={filterStatusOptions}
	onSearch={resetAndFetchSources}
/>
<button onclick={resetAndFetchSources} disabled={loading}>Refresh Sources</button>
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
			<button onclick={fetchSources} disabled={loading}>Load More</button>
		{/if}
	{:else}
		<p>No sources found.</p>
	{/if}
{:else}
	<p>No sources found.</p>
{/if}
<PageSizeSelector bind:pageSize />
