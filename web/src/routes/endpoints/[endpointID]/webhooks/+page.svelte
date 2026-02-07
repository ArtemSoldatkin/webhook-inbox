<script lang="ts">
	import { page } from "$app/state";
	import type { Webhook } from "$lib/types";
	import { onMount } from "svelte";

    const {endpointID} = page.params
    let data: Webhook[] | null = null
    let loading = false
    let error: string | null = null

    async function fetchWebhooks() {
		loading = true;
		error = null;
		try {
			const res = await fetch(`/api/webhooks?endpointID=${endpointID}`);
			if (!res.ok) {
				throw new Error('Failed to fetch data');
			}
			data = await res.json();
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching webhooks:', err);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		if (endpointID) {
			fetchWebhooks();
		} else {
			error = 'No endpoint ID provided in URL';
		}
	});
</script>
{#if loading}
	<p>Loading...</p>
{:else if error}
	<p>Error: {error}</p>
{:else if data}
	{#if data.length === 0}
		<p>No webhooks found for this user.</p>
	{/if}
	{#each data as webhook}
		<div>
			<h2>{webhook.Name}</h2>
			<p>{webhook.Description}</p>
			<p>{webhook.PublicKey}</p>
			<p>Active: {webhook.IsActive ? 'Yes' : 'No'}</p>
			<p>Created At: {new Date(webhook.CreatedAt).toLocaleString()}</p>
			<p>Updated At: {new Date(webhook.UpdatedAt).toLocaleString()}</p>
		</div>
	{/each}
{/if}
