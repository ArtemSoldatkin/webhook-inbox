<script lang="ts">
	import InputMap from '$lib/components/InputMap.svelte';
	import env from '$lib/env';
	import type { SourceDTO } from '$lib/types';

	type NewSource = Omit<
		SourceDTO,
		| 'id'
		| 'public_id'
		| 'ingress_url'
		| 'status'
		| 'status_reason'
		| 'created_at'
		| 'updated_at'
		| 'disable_at'
	>;

	let data = newData();
	let loading = false;
	let error: string | null = null;

	function newData(): NewSource {
		return { egress_url: '', static_headers: {}, description: '' };
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
		return validateEgressUrl(data.egress_url);
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
			const parsedUrl = new URL(url);
			if (
				(parsedUrl.protocol !== 'http:' && parsedUrl.protocol !== 'https:') ||
				parsedUrl.href.length > 2048
			) {
				return false;
			}
			if (env.VITE_ENV === 'dev') return true;
			return (
				/^https?:\/\//.test(parsedUrl.href) &&
				!/^https?:\/\/(localhost|127\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})|0\.0\.0\.0|\[?::1\]?)(\/|:|$)/.test(
					parsedUrl.href
				) &&
				!/^https?:\/\/\[\:\:ffff\:127\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\]/.test(parsedUrl.href) &&
				!/^https?:\/\/10\./.test(parsedUrl.href) &&
				!/^https?:\/\/192\.168\./.test(parsedUrl.href) &&
				!/^https?:\/\/172\.(1[6-9]|2[0-9]|3[0-1])\./.test(parsedUrl.href) &&
				!/^https?:\/\/169\.254\.169\.254(\/|:|$)/.test(parsedUrl.href) &&
				!/^https?:\/\/\[::ffff:0\.0\.0\.0\]/.test(parsedUrl.href) &&
				!/^https?:\/\/localhost\.(\/|:|$)/.test(parsedUrl.href)
			);
		} catch {
			return false;
		}
	}

	let egressError: string | null = null;
	$: if (data.egress_url.trim() !== '' && !validateEgressUrl(data.egress_url)) {
		egressError = 'Valid Egress URL is required';
	} else {
		egressError = null;
	}
</script>

<form on:submit|preventDefault={handleSubmit}>
	<label
		>Egress URL
		<input
			type="text"
			bind:value={data.egress_url}
			placeholder="https://example.com/egress"
			required
			disabled={loading}
		/>
	</label>
	{#if egressError}
		<p class="error">{egressError}</p>
	{/if}
	<label
		>Static Headers
		<InputMap bind:json={data.static_headers} disabled={loading} />
	</label>
	<label
		>Description
		<textarea bind:value={data.description} placeholder="Optional description"></textarea>
	</label>
	<button type="submit" disabled={loading || Boolean(egressError)}>Create New Source</button>
	{#if error}
		<p class="error">Error: {error}</p>
	{/if}
</form>
