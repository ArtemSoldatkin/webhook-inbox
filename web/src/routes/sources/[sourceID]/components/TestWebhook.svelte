<script lang="ts">
	import InputMap from '$lib/components/InputMap.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import type { ContentType } from '$lib/types';
	import BodyInput from './BodyInput.svelte';

	type Props = {
		/** Public source identifier used by the ingest endpoint. */
		publicID: string;

		/** Static headers automatically sent with the test request. */
		staticHeaders?: Record<string, string>;
	};

	/** Data required to send a test webhook request. */
	let { publicID, staticHeaders }: Props = $props();

	/** Supported HTTP methods for the test form. */
	type HTTPMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

	/** Selected HTTP method for the test request. */
	let method = $state<HTTPMethod>('GET');

	/** Additional request headers entered by the user. */
	let headers = $state<Record<string, string>>({});

	/** Query parameters appended to the test request. */
	let queryParams = $state<Record<string, string>>({});

	/** Text-based body value shared by non-form-data editors. */
	let textBody = $state('');

	/** Multipart body payload built by the form-data editor. */
	let formDataBody = $state(new FormData());

	/** Selected content type for the request body. */
	let contentType = $state<ContentType>('application/json');

	/** Tracks whether a test request is in flight. */
	let loading = $state(false);

	/** Holds the latest webhook test error. */
	let error = $state<string | null>(null);

	/** Normalized request body for the current content type. */
	let body = $derived(
		contentType === 'multipart/form-data'
			? formDataBody
			: contentType === 'application/octet-stream' && textBody
				? base64ToUint8Array(textBody)
				: textBody
	);

	/** Indicates whether the selected method accepts a body. */
	let isBodyAllowed = $derived(method !== 'GET');

	/** Indicates whether the form should send a Content-Type header. */
	let isContentTypeAllowed = $derived(isBodyAllowed && contentType !== 'multipart/form-data');

	/**
	 * Decodes a base64 string into request bytes.
	 *
	 * @param base64 - Base64-encoded request body.
	 * @returns Binary body as a byte array.
	 */
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

	/**
	 * Sends the configured test webhook request.
	 *
	 * @param event - Form submission event.
	 */
	async function testWebhook(event: SubmitEvent): Promise<void> {
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
	<Button type="submit" disabled={loading}>Test Webhook</Button>
</form>
