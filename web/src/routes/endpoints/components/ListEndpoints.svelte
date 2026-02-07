<script lang="ts">
	import type { Endpoint } from '$lib/types';

	let userID: string = '';
	let data: Endpoint[] | null = null;
	let loading = false;
	let error: string | null = null;

	async function fetchEndpoints() {
		loading = true;
		error = null;
		try {
			const res = await fetch(`api/endpoints?user_id=${userID}`);
			if (!res.ok) {
				throw new Error('Failed to fetch data');
			}
			data = await res.json();
			userID = '';
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching endpoints:', err);
		} finally {
			loading = false;
		}
	}
</script>

<h2>Upload endpoints</h2>
<form on:submit|preventDefault={fetchEndpoints}>
	<label for="userID">User ID:</label>
	<input
		type="text"
		name="userID"
		placeholder="Enter user ID"
		required
		bind:value={userID}
		disabled={loading}
	/>
	<button type="submit" disabled={loading}>Fetch Data</button>
</form>
{#if loading}
	<p>Loading...</p>
{:else if error}
	<p>Error: {error}</p>
{:else if data}
	{#if data.length === 0}
		<p>No endpoints found for this user.</p>
	{/if}
	{#each data as endpoint}
		<div>
			<h2>{endpoint.Name}</h2>
			<p>{endpoint.Description}</p>
			<p>{endpoint.Url}</p>
			{#each Object.entries(endpoint.Headers) as [headerKey, headerValue]}
				<p>{headerKey}: {headerValue}</p>
			{/each}
			<p>Active: {endpoint.IsActive ? 'Yes' : 'No'}</p>
			<p>Created At: {new Date(endpoint.CreatedAt).toLocaleString()}</p>
		</div>
	{/each}
{/if}
