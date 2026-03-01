<script lang="ts">
	export let url: string;
	export let defaultPageSize: number = 20;

	let data: any[] = [];
	let loading = false;
	let error: string | null = null;

	let pageSize: number = defaultPageSize;
	let nextCursor: string | null = null;
	let hasNext: boolean = false;

	async function fetchData() {
		if (!url) {
			error = 'URL is required to fetch data.';
			return;
		}
		loading = true;
		error = null;
		try {
			const params = new URLSearchParams({
				limit: pageSize.toString(),
				cursor: nextCursor ?? ''
			});
			const response = await fetch(`${url}?${params.toString()}`);
			if (!response.ok) {
				throw new Error(`Failed to fetch data: ${response.statusText}`);
			}
			const rawData = await response.json();
			data.push(...rawData.Data);
			nextCursor = rawData.NextCursor;
			hasNext = rawData.HasNext;
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching data:', err);
		} finally {
			loading = false;
		}
	}

	$: if (url) {
		fetchData();
	}

	$: if (defaultPageSize !== pageSize) {
		pageSize = defaultPageSize;
		data = [];
		nextCursor = null;
		hasNext = false;
		fetchData();
	}
</script>

<section>
	{#if loading}
		<p>Loading...</p>
	{:else if error}
		<p class="error">Error: {error}</p>
	{:else if data.length > 0}
		<ul>
			{#each data as item}
				<slot name="item" {item} />
			{/each}
		</ul>
	{/if}

	<select bind:value={pageSize}>
		<option value="10">10</option>
		<option value="20">20</option>
		<option value="50">50</option>
		<option value="100">100</option>
	</select>

	{#if hasNext}
		<button on:click={fetchData} disabled={loading}>Load More</button>
	{/if}
</section>
