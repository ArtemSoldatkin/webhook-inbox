<script lang="ts">
	import { fetchPaginatedData } from '$lib/api';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import { parseSourceDTO } from '$lib/dtoParsers';
	import type { SourceDTO } from '$lib/types';
	import { onMount } from 'svelte';

	let data: SourceDTO[] = [];
	let loading = false;
	let error: string | null = null;

	let pageSize: number = 20;
	let nextCursor: string | null = null;
	let hasNext: boolean = false;

	async function fetchSources() {
		loading = true;
		error = null;
		try {
			const result = await fetchPaginatedData('/api/sources', pageSize, nextCursor);
			data = [...data, ...result.data.map(parseSourceDTO)];
			nextCursor = result.nextCursor;
			hasNext = result.hasNext;
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

	onMount(() => {
		fetchSources();
	});

	$: if (pageSize) {
		resetAndFetchSources();
	}
</script>

<button on:click={resetAndFetchSources} disabled={loading}>Refresh Sources</button>
{#if loading}
	<p>Loading sources...</p>
{:else if error}
	<p class="error">Error: {error}</p>
{:else if data}
	{#if data.length > 0}
		<ul>
			{#each data as source}
				<li>
					<section>
						<h2><a href={`/sources/${source.ID}`}>{source.ID}</a></h2>
						<p>{source.Description}</p>
						<p>{source.IngressUrl}</p>
						<p>{source.EgressUrl}</p>
						<p>Static headers:</p>
						{#each Object.entries(source.StaticHeaders) as [key, value]}
							<p>{key}: {value}</p>
						{/each}
						<p>{source.Status}</p>
						<p>{source.StatusReason}</p>
						<p>Created at: {new Date(source.CreatedAt).toLocaleString()}</p>
						<p>Updated at: {new Date(source.UpdatedAt).toLocaleString()}</p>
						<p>
							Disabled at: {source.DisableAt ? new Date(source.DisableAt).toLocaleString() : 'N/A'}
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
