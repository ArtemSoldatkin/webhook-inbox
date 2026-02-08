<script lang="ts">
	import { userID } from '$lib/stores/global';
	import type { Endpoint } from '$lib/types';
	import { onMount } from 'svelte';

	let data: Endpoint[] | null = null;
	let loading = false;
	let error: string | null = null;

	async function fetchEndpoints() {
		loading = true;
		error = null;
		try {
			const res = await fetch(`/api/endpoints?userID=${$userID}`);
			if (!res.ok) {
				throw new Error('Failed to fetch data');
			}
			data = await res.json();
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching endpoints:', err);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchEndpoints();
	});
</script>

<h2>Endpoints list</h2>
<nav>
	<ul>
		<li><a href="/endpoints">Endpoints</a></li>
	</ul>
</nav>
{#if loading}
	<p>Loading...</p>
{:else if error}
	<p>Error: {error}</p>
{:else if data}
	{#if data.length === 0}
		<p>No endpoints found for this user.</p>
	{/if}
	<ul>
		{#each data as endpoint}
			<li>
				<section>
					<h3><a href={`/endpoints/${endpoint.ID}`}>{endpoint.Name}</a></h3>
					<p>{endpoint.Description}</p>
					<p>{endpoint.Url}</p>
					{#each Object.entries(endpoint.Headers) as [headerKey, headerValue]}
						<p>{headerKey}: {headerValue}</p>
					{/each}
					<p>Active: {endpoint.IsActive ? 'Yes' : 'No'}</p>
					<p>Created At: {new Date(endpoint.CreatedAt).toLocaleString()}</p>
				</section>
			</li>
		{/each}
	</ul>
{/if}
