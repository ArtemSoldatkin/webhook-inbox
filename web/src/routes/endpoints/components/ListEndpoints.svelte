<script lang="ts">
	interface Endpoint {
		ID: number;
		UserID: number;
		Url: string;
		Name: string;
		Description: string;
		Headers: string[];
		IsActive: boolean;
		CreatedAt: Date;
	}

	let userId: string = '';
	let data: Endpoint[] | null = null;
	let loading = false;
	let error: string | null = null;

	async function fetchEndpoints() {
		loading = true;
		error = null;
		try {
			const res = await fetch(`api/endpoints?user_id=${userId}`);
			if (!res.ok) {
				throw new Error('Failed to fetch data');
			}
			data = await res.json();
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
		} finally {
			loading = false;
		}
	}
</script>

<h2>Upload endpoints</h2>
<form on:submit|preventDefault={fetchEndpoints}>
	<label for="userId">User ID:</label>
	<input
		type="text"
		name="userId"
		placeholder="Enter user ID"
		required
		bind:value={userId}
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
			<p>{endpoint.Url}</p>
			<p>{endpoint.Description}</p>
			<p>Active: {endpoint.IsActive ? 'Yes' : 'No'}</p>
			<p>Created At: {new Date(endpoint.CreatedAt).toLocaleString()}</p>
		</div>
	{/each}
{/if}
