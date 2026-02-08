<script lang="ts">
	import { page } from '$app/state';
	import type { Webhook } from '$lib/types';
	import { onMount } from 'svelte';

	const { endpointID } = page.params;
	let data: Webhook[] | null = null;
	let loading = false;
	let error: string | null = null;

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

<nav>
	<ul>
		<li><a href={`/endpoints/${endpointID}/webhooks`}>Webhooks</a></li>
	</ul>
</nav>
<h2>Webhooks for endpoint {endpointID}</h2>
{#if loading}
	<p>Loading...</p>
{:else if error}
	<p>Error: {error}</p>
{:else if data}
	{#if data.length === 0}
		<p>No webhooks found for this user.</p>
	{/if}
	<ul>
		{#each data as webhook}
			<li>
				<section>
					<h3><a href={`/endpoints/${endpointID}/webhooks/${webhook.ID}`}>{webhook.Name}</a></h3>
					<p>{webhook.Description}</p>
					<p>{webhook.PublicKey}</p>
					<p>Active: {webhook.IsActive ? 'Yes' : 'No'}</p>
					<p>Created At: {new Date(webhook.CreatedAt).toLocaleString()}</p>
					<p>Updated At: {new Date(webhook.UpdatedAt).toLocaleString()}</p>
				</section>
			</li>
		{/each}
	</ul>
{/if}
