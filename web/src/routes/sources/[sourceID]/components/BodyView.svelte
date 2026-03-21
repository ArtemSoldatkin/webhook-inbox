<script lang="ts">
	import type { ContentType } from '$lib/types';
	import ByteBodyView from './BodyTypeViews/ByteBodyView.svelte';
	import FormUrlEncodedBodyView from './BodyTypeViews/FormUrlEncodedBodyView.svelte';
	import JSONBodyView from './BodyTypeViews/JSONBodyView.svelte';
	import PlainTextBodyView from './BodyTypeViews/PlainTextBodyView.svelte';
	import XMLBodyView from './BodyTypeViews/XMLBodyView.svelte';

	type Props = {
		body?: string;
		contentType?: ContentType;
	};

	type ParsedBody = {
		content: string;
		error: string | null;
	};

	let { body, contentType }: Props = $props();
	const parsedBody = $derived(parseBody(body, contentType));

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
