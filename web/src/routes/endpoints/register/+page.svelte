<script lang="ts">
	import InputJSON from '$lib/components/InputJSON.svelte';
	import type { Endpoint } from '$lib/types';

	let userID: string = '';
	let url: string = '';
	let name: string = '';
	let description: string = '';
	let headers: Record<string, string> = {};

	let data: Endpoint | null = null;
	let loading = false;
	let error: string | null = null;

	async function registerEndpoint() {
		loading = true;
		try {
			const res = await fetch('api/endpoints', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					userID: parseInt(userID, 10),
					url,
					name,
					description,
					headers
				})
			});
			if (!res.ok) {
				throw new Error('Failed to register endpoint');
			}
			data = await res.json();
			userID = '';
			url = '';
			name = '';
			description = '';
			headers = {};
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error registering endpoints:', err);
		} finally {
			loading = false;
		}
	}
</script>

<h2>Register a new endpoint</h2>
<form on:submit|preventDefault={registerEndpoint}>
	<label for="userID">User ID:</label>
	<input
		type="text"
		name="userID"
		placeholder="Enter user ID"
		required
		bind:value={userID}
		disabled={loading}
	/>
	<label for="url">URL:</label>
	<input
		type="text"
		name="url"
		placeholder="Enter endpoint URL"
		required
		bind:value={url}
		disabled={loading}
	/>
	<label for="name">Name:</label>
	<input
		type="text"
		name="name"
		placeholder="Enter endpoint name"
		required
		bind:value={name}
		disabled={loading}
	/>
	<label for="description">Description:</label>
	<textarea
		name="description"
		placeholder="Enter endpoint description"
		bind:value={description}
		disabled={loading}
	></textarea>
	<label
		>Headers (JSON format):
		<InputJSON bind:json={headers} disabled={loading} />
	</label>

	<button type="submit" disabled={loading}>Register Endpoint</button>
</form>
{#if loading}
	<p>Loading...</p>
{:else if error}
	<p>Error: {error}</p>
{:else if data}
	<div>
		<h2>{data.Name}</h2>
		<p>{data.Description}</p>
		<p>{data.Url}</p>
		{#each Object.entries(data.Headers) as [headerKey, headerValue]}
			<p>{headerKey}: {headerValue}</p>
		{/each}
		<p>Active: {data.IsActive ? 'Yes' : 'No'}</p>
		<p>Created At: {new Date(data.CreatedAt).toLocaleString()}</p>
	</div>
{/if}
