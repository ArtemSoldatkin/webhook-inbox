<script lang="ts">
	import InputMap from '$lib/components/InputMap.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Select from '$lib/components/ui/Select.svelte';
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

<form onsubmit={testWebhook} class="flex flex-col gap-6">
	<div class="grid gap-6 lg:grid-cols-[minmax(0,0.36fr)_minmax(0,0.64fr)]">
		<div class="flex flex-col gap-2">
			<label class="text-sm font-medium text-fg">
				HTTP Method
				<Select
					bind:value={method}
					disabled={loading}
					options={[
						{ value: 'GET', label: 'GET' },
						{ value: 'POST', label: 'POST' },
						{ value: 'PUT', label: 'PUT' },
						{ value: 'PATCH', label: 'PATCH' },
						{ value: 'DELETE', label: 'DELETE' }
					]}
					class="mt-2 rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
				/>
			</label>
		</div>
		<Alert>
			Requests are sent to the source ingest endpoint using the selected method, merged static
			headers, and any overrides you add below.
		</Alert>
	</div>

	<div class="flex flex-col gap-3 rounded-lg border border-border-muted bg-elevated p-4">
		<div>
			<p class="text-sm font-medium text-fg">Headers</p>
			<p class="mt-1 text-sm text-muted">Optional headers included with the test request.</p>
		</div>
		<InputMap bind:map={headers} disabled={loading} />
	</div>

	<div class="flex flex-col gap-3 rounded-lg border border-border-muted bg-elevated p-4">
		<div>
			<p class="text-sm font-medium text-fg">Query parameters</p>
			<p class="mt-1 text-sm text-muted">Optional query string values appended to the request URL.</p>
		</div>
		<InputMap bind:map={queryParams} disabled={loading} />
	</div>

	{#if isBodyAllowed}
		<div class="flex flex-col gap-3 rounded-lg border border-border-muted bg-elevated p-4">
			<div>
				<p class="text-sm font-medium text-fg">Request body</p>
				<p class="mt-1 text-sm text-muted">Choose a content type and compose the request payload.</p>
			</div>
			<BodyInput bind:textBody bind:formDataBody bind:contentType />
		</div>
	{/if}

	{#if error}
		<Alert variant="error" title="Error" class="bg-surface">{error}</Alert>
	{/if}

	<div class="flex justify-end border-t border-border-muted pt-6">
		<Button type="submit" disabled={loading}>
			{loading ? 'Sending Test Webhook...' : 'Test Webhook'}
		</Button>
	</div>
</form>
