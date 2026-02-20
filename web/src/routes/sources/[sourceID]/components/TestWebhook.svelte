<script lang="ts">
	import InputMap from '$lib/components/InputMap.svelte';
	import type { ContentType } from '$lib/types';
	import BodyInput from './BodyInput.svelte';

	export let publicID: string;
	export let staticHeaders: Record<string, string> = {};

	type HTTPMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

	let method: HTTPMethod = 'GET';
	let headers: Record<string, string> = {};
	let queryParams: Record<string, string> = {};
	let body: string = '';
	let contentType: ContentType = 'application/json';

	let loading = false;
	let error: string | null = null;

	$: isBodyAllowed = method !== 'GET';

	async function testWebhook() {
		loading = true;
		error = null;
		try {
			const webhookURL = `/api/ingest/${publicID}`;
			const params = new URLSearchParams(queryParams).toString();
			const urlWithParams = params ? `${webhookURL}?${params}` : webhookURL;
			const response = await fetch(urlWithParams, {
				method,
				headers: {
					...staticHeaders,
					...headers,
					...(isBodyAllowed && { 'Content-Type': contentType })
				},
				...(isBodyAllowed && { body })
			});
			if (!response.ok) {
				throw new Error(`Failed to test webhook: ${response.statusText}`);
			}
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error testing webhook:', err);
		} finally {
			loading = false;
		}
	}
</script>

<form on:submit|preventDefault={testWebhook}>
	<label
		>HTTP Method:
		<select bind:value={method} disabled={loading}>
			<option value="GET">GET</option>
			<option value="POST">POST</option>
			<option value="PUT">PUT</option>
			<option value="PATCH">PATCH</option>
			<option value="DELETE">DELETE</option>
		</select>
	</label>
	<label>
		Headers (optional):
		<InputMap bind:json={headers} disabled={loading} />
	</label>
	<label>
		Query Parameters (optional):
		<InputMap bind:json={queryParams} disabled={loading} />
	</label>
	{#if isBodyAllowed}<label>
			Body (optional):
			<BodyInput bind:body bind:contentType />
		</label>
	{/if}
	{#if error}
		<div class="error">{error}</div>
	{/if}
	<button type="submit" disabled={loading}>Test Webhook</button>
</form>
