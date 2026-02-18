<script lang="ts">
	import InputMap from '$lib/components/InputMap.svelte';
	import env from '$lib/env';
	import type { SourceDTO } from '$lib/types';

	type NewSource = Omit<
		SourceDTO,
		| 'ID'
		| 'PublicID'
		| 'IngressUrl'
		| 'Status'
		| 'StatusReason'
		| 'CreatedAt'
		| 'UpdatedAt'
		| 'DisableAt'
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
		return validateEgressUrl(data.EgressUrl);
	}

	function handleSubmit() {
		const isValid = validateInput();
		if (!isValid) {
			console.warn('Invalid input, cannot create source');
			return;
		}
		createSource();
	}

	function validateEgressUrl(url: string): boolean {
		try {
			new URL(url);
			if (env.VITE_ENV === 'dev') {
				return true;
			}
			// TODO add additional checks to ensure the URL is valid for egress (e.g., not localhost, etc.)
			return true;
		} catch {
			return false;
		}
	}

	let egressError: string | null = null;
	$: if (data.EgressUrl.trim() !== '' && !validateEgressUrl(data.EgressUrl)) {
		egressError = 'Egress URL is required';
	} else {
		egressError = null;
	}
</script>

<form on:submit|preventDefault={handleSubmit}>
	<label
		>Egress URL
		<input
			type="text"
			bind:value={data.EgressUrl}
			placeholder="https://example.com/egress"
			required
			disabled={loading}
		/>
	</label>
	{#if egressError}
		<p class="error">Egress URL Error: {egressError}</p>
	{/if}
	<label
		>Static Headers
		<InputMap bind:json={data.StaticHeaders} disabled={loading} />
	</label>
	<label
		>Description
		<textarea bind:value={data.Description} placeholder="Optional description"></textarea>
	</label>
	<button type="submit" disabled={loading || Boolean(egressError)}>Create New Source</button>
	{#if error}
		<p class="error">Error: {error}</p>
	{/if}
</form>
