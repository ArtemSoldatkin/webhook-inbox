<script lang="ts">
	import type { Source } from '$lib/types';
	import { onMount } from 'svelte';

	let data: Source[] | null = null;
	let loading = false;
	let error: string | null = null;

	async function fetchSources() {
		loading = true;
		error = null;
		try {
			const response = await fetch('/api/sources');
			if (!response.ok) {
				throw new Error(`Failed to fetch sources: ${response.statusText}`);
			}
			data = await response.json();
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching sources:', err);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchSources();
	});
</script>

<button on:click={fetchSources}>Refresh Sources</button>
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
						<h2>{source.ID}</h2>
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
	{:else}
		<p>No sources found.</p>
	{/if}
{:else}
	<p>No sources found.</p>
{/if}
