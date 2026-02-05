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

	let data: Endpoint[] | null = null;
	let loading = false;
	let error: string | null = null;

	async function fetchData() {
		loading = true;
		error = null;
		try {
			const res = await fetch('/api/endpoints?user_id=2');
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

{#if loading}
	<p>Loading...</p>
{:else if error}
	<p>Error: {error}</p>
{:else if data}
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

<button on:click={fetchData}>Fetch Data</button>
