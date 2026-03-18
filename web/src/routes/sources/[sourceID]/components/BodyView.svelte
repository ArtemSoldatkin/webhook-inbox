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

	let { body, contentType }: Props = $props();

	const parsedBody = $derived(body ? atob(body) : '');
</script>

<section>
	<h3>Request body</h3>
	{#if !body}
		<p>No body content</p>
	{:else if !contentType}
		<p>Content type unknown, cannot display body</p>
	{:else if contentType.startsWith('application/json')}
		<JSONBodyView body={parsedBody} />
	{:else if contentType.startsWith('application/x-www-form-urlencoded')}
		<FormUrlEncodedBodyView body={parsedBody} />
	{:else if contentType.startsWith('application/xml')}
		<XMLBodyView body={parsedBody} />
	{:else if contentType.startsWith('text/plain')}
		<PlainTextBodyView body={parsedBody} />
	{:else}
		<ByteBodyView body={parsedBody} {contentType} />
	{/if}
</section>
