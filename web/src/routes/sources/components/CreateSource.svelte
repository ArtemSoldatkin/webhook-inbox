<script lang="ts">
	import InputMap from '$lib/components/InputMap.svelte';
	import type { SourceDTO } from '$lib/types';

	type NewSource = Omit<
		SourceDTO,
		'ID' | 'IngressUrl' | 'Status' | 'StatusReason' | 'CreatedAt' | 'UpdatedAt' | 'DisableAt'
	>;

	let data = newData();
	let loading = false;
	let error: string | null = null;

	function newData(): NewSource {
		return { EgressUrl: '', StaticHeaders: {}, Description: '' };
	}

	async function createSource() {
		loading = true;
		error = null;
		try {
			const response = await fetch('/api/sources', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(data)
			});
			if (!response.ok) {
				throw new Error(`Failed to create source: ${response.statusText}`);
			}
			const newSource = await response.json();
			data = newData(); // Reset form after successful creation
			console.log('Created new source:', newSource);
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error creating source:', err);
		} finally {
			loading = false;
		}
	}

	function validateInput() {
		if (data.EgressUrl.trim() === '') {
			return false;
		}
		return true;
	}

	function handleSubmit() {
		const isValid = validateInput();
		if (!isValid) {
			console.warn('Invalid input, cannot create source');
			return;
		}
		createSource();
	}
</script>

<form on:submit|preventDefault={handleSubmit}>
	<label for="egress-url">Egress URL</label>
	<input
		id="egress-url"
		type="text"
		bind:value={data.EgressUrl}
		placeholder="https://example.com/egress"
		required
		disabled={loading}
	/>
	<label
		>Static Headers
		<InputMap bind:json={data.StaticHeaders} disabled={loading} />
	</label>
	<label for="description">Description</label>
	<textarea id="description" bind:value={data.Description} placeholder="Optional description"
	></textarea>
	<button type="submit" disabled={loading}>Create New Source</button>
	{#if error}
		<p class="error">Error: {error}</p>
	{/if}
</form>
