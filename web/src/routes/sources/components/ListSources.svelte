<script lang="ts">
	import { resolve } from '$app/paths';
	import { fetchPaginatedData } from '$lib/api';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import { parseSourceDTO } from '$lib/dtoParsers';
	import type { SourceDTO } from '$lib/types';

	let data: SourceDTO[] = [];
	let loading = false;
	let error: string | null = null;

	let pageSize: number = 20;
	let nextCursor: string | null = null;
	let hasNext: boolean = false;

	let searchQuery: string = '';

	let filterStatus: string = '';
	const filterStatusOptions = ['active', 'paused', 'quarantined', 'disabled'];

	async function fetchSources() {
		loading = true;
		error = null;
		try {
			const result = await fetchPaginatedData('/api/sources', pageSize, nextCursor, {
				search: searchQuery,
				filter_status: filterStatus
			});
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

	$: if (pageSize && filterStatus) {
		resetAndFetchSources();
	}
</script>

<FilterBar
	bind:searchQuery
	bind:filter={filterStatus}
	filterOptions={filterStatusOptions}
	onSearch={resetAndFetchSources}
/>
<button on:click={resetAndFetchSources} disabled={loading}>Refresh Sources</button>
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
			<button on:click={fetchSources} disabled={loading}>Load More</button>
		{/if}
	{:else}
		<p>No sources found.</p>
	{/if}
{:else}
	<p>No sources found.</p>
{/if}
<PageSizeSelector bind:pageSize />
