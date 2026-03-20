<script lang="ts">
	import InputMap from '$lib/components/InputMap.svelte';
	import type { ContentType } from '$lib/types';
	import BodyInput from './BodyInput.svelte';

	type Props = {
		publicID: string;
		staticHeaders?: Record<string, string>;
	};

	let { publicID, staticHeaders }: Props = $props();

	type HTTPMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

	let method = $state<HTTPMethod>('GET');
	let headers = $state<Record<string, string>>({});
	let queryParams = $state<Record<string, string>>({});
	let textBody = $state('');
	let formDataBody = $state(new FormData());

	let contentType = $state<ContentType>('application/json');

	let loading = $state(false);
	let error = $state<string | null>(null);

	let body = $derived(
		contentType === 'multipart/form-data'
			? formDataBody
			: contentType === 'application/octet-stream' && textBody
				? base64ToUint8Array(textBody)
				: textBody
	);

	let isBodyAllowed = $derived(method !== 'GET');
	let isContentTypeAllowed = $derived(isBodyAllowed && contentType !== 'multipart/form-data');

	function base64ToUint8Array(base64: string): Uint8Array<ArrayBuffer> {
		if (!base64) return new Uint8Array(0);
		const binaryString = atob(base64);
		const buffer = new ArrayBuffer(binaryString.length);
		const bytes = new Uint8Array(buffer);
		for (let i = 0; i < binaryString.length; i++) {
			bytes[i] = binaryString.charCodeAt(i) & 0xff;
		}
		return bytes;
	}

	async function testWebhook(event: SubmitEvent) {
		event.preventDefault();
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
					...(isContentTypeAllowed && { 'Content-Type': contentType })
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

<form onsubmit={testWebhook}>
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
		<InputMap bind:map={headers} disabled={loading} />
	</label>
	<label>
		Query Parameters (optional):
		<InputMap bind:map={queryParams} disabled={loading} />
	</label>
	{#if isBodyAllowed}<label>
			Body (optional):
			<BodyInput bind:textBody bind:formDataBody bind:contentType />
		</label>
	{/if}
	{#if error}
		<div class="error">{error}</div>
	{/if}
	<button type="submit" disabled={loading}>Test Webhook</button>
</form>
