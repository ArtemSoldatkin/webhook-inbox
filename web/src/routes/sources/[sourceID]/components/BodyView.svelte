<script lang="ts">
	import type { ContentType } from '$lib/types';
	import ByteBodyView from './body/views/ByteBodyView.svelte';
	import FormUrlEncodedBodyView from './body/views/FormUrlEncodedBodyView.svelte';
	import JSONBodyView from './body/views/JSONBodyView.svelte';
	import PlainTextBodyView from './body/views/PlainTextBodyView.svelte';
	import XMLBodyView from './body/views/XMLBodyView.svelte';

	type Props = {
		/** Raw stored body value to display. */
		body?: string;

		/** Content type associated with the stored body. */
		contentType?: ContentType;
	};

	/** Normalized body content ready for display. */
	type ParsedBody = {
		content: string;
		error: string | null;
	};

	let { body, contentType }: Props = $props();

	/** Parsed version of the current body payload. */
	const parsedBody = $derived(parseBody(body, contentType));

	/**
	 * Decodes the stored body into displayable text or binary content.
	 *
	 * @param body - Raw stored body value.
	 * @param contentType - Content type associated with the body.
	 * @returns Parsed body content and any display error.
	 */
	function parseBody(body?: string, contentType?: ContentType): ParsedBody {
		if (body === undefined)
			return {
				content: '',
				error: 'No body provided, cannot display content'
			};

		if (body === '')
			return {
				content: '',
				error: 'Body content is empty'
			};

		try {
			const binary = atob(body);
			if (
				!contentType ||
				contentType.startsWith('multipart/form-data') ||
				contentType.startsWith('application/octet-stream')
			) {
				return {
					content: binary,
					error: null
				};
			}

			const bytes = new Uint8Array(binary.length);

			for (let i = 0; i < binary.length; i++) {
				bytes[i] = binary.charCodeAt(i);
			}

			return {
				content: new TextDecoder('utf-8', { fatal: false }).decode(bytes),
				error: null
			};
		} catch (err) {
			console.error('Failed to parse body as Base64 string', err);

			return {
				content: body,
				error: 'Failed to parse body as Base64 string, falling back to raw content'
			};
		}
	}
</script>

<section>
	<h3>Request body</h3>
	{#if parsedBody.error && parsedBody.content}
		<p>{parsedBody.error}</p>
		<p>Original body: {parsedBody.content}</p>
	{:else if parsedBody.error}
		<p>{parsedBody.error}</p>
	{:else if !contentType}
		<p>Content type unknown, cannot display body</p>
	{:else if contentType.startsWith('application/json')}
		<JSONBodyView body={parsedBody.content} />
	{:else if contentType.startsWith('application/x-www-form-urlencoded')}
		<FormUrlEncodedBodyView body={parsedBody.content} />
	{:else if contentType.startsWith('application/xml')}
		<XMLBodyView body={parsedBody.content} />
	{:else if contentType.startsWith('text/plain')}
		<PlainTextBodyView body={parsedBody.content} />
	{:else}
		<ByteBodyView body={parsedBody.content} {contentType} />
	{/if}
</section>
