<script lang="ts">
	import { parseSourceDTO } from '$lib/dtoParsers';
	import type { SourceDTO } from '$lib/types';
	import ListEvents from './ListEvents.svelte';
	import TestWebhook from './TestWebhook.svelte';

	export let sourceID: string;

	let data: SourceDTO | null = null;
	let loading = false;
	let error: string | null = null;

	async function fetchSource() {
		loading = true;
		error = null;
		try {
			const response = await fetch(`/api/sources/${sourceID}`);
			if (!response.ok) {
				throw new Error(`Failed to fetch source: ${response.statusText}`);
			}
			const rawData = await response.json();
			data = parseSourceDTO(rawData);
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching source:', err);
		} finally {
			loading = false;
		}
	}

	$: if (sourceID) {
		fetchSource();
	}
</script>

{#if loading}
	<p>Loading source data...</p>
{:else if error}
	<p class="error">Error: {error}</p>
{:else if data}
	<section>
		<h2>{data.ID}</h2>
		<p>{data.Description}</p>
		<p>{data.IngressUrl}</p>
		<p>{data.EgressUrl}</p>
		<p>Static headers:</p>
		{#each Object.entries(data.StaticHeaders) as [key, value]}
			<p>{key}: {value}</p>
		{/each}
		<p>{data.Status}</p>
		<p>{data.StatusReason}</p>
		<p>Created at: {new Date(data.CreatedAt).toLocaleString()}</p>
		<p>Updated at: {new Date(data.UpdatedAt).toLocaleString()}</p>
		<p>
			Disabled at: {data.DisableAt ? new Date(data.DisableAt).toLocaleString() : 'N/A'}
		</p>
	</section>
	<section>
		<h3>Test Webhook</h3>
		<TestWebhook publicID={data.PublicID} staticHeaders={data.StaticHeaders} />
		<ListEvents {sourceID} />
	</section>
{/if}
