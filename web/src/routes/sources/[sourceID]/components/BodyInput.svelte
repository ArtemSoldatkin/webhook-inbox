<script lang="ts">
	import type { ContentType } from '$lib/types';
	import ByteBodyInput from './body/inputs/ByteBodyInput.svelte';
	import FormDataBodyInput from './body/inputs/FormDataBodyInput.svelte';
	import FormUrlEncodedBodyInput from './body/inputs/FormUrlEncodedBodyInput.svelte';
	import JSONBodyInput from './body/inputs/JSONBodyInput.svelte';
	import PlainTextBodyInput from './body/inputs/PlainTextBodyInput.svelte';
	import XMLBodyInput from './body/inputs/XMLBodyInput.svelte';

	type Props = {
		/** Bound text body used by text-based editors. */
		textBody: string;

		/** Bound multipart body used by the form-data editor. */
		formDataBody: FormData;

		/** Selected content type for the active body editor. */
		contentType: ContentType;
	};

	/** Bound body inputs shared with the webhook form. */
	let {
		textBody = $bindable(),
		formDataBody = $bindable(),
		contentType = $bindable()
	}: Props = $props();

	/** Validation error returned by the active body editor. */
	let error = $state<string | null>(null);

	$effect(() => {
		if (!contentType) return;
		textBody = '';
		formDataBody = new FormData();
		error = null;
	});
</script>

<section>
	<select bind:value={contentType}>
		<option value="application/json">JSON</option>
		<option value="application/x-www-form-urlencoded">Form URL Encoded</option>
		<option value="multipart/form-data">Multipart Form Data</option>
		<option value="text/plain">Plain Text</option>
		<option value="application/xml">XML</option>
		<option value="application/octet-stream">Binary Data</option>
	</select>

	{#if contentType === 'application/json'}
		<JSONBodyInput bind:body={textBody} bind:error />
	{:else if contentType === 'application/x-www-form-urlencoded'}
		<FormUrlEncodedBodyInput bind:body={textBody} bind:error />
	{:else if contentType === 'multipart/form-data'}
		<FormDataBodyInput bind:body={formDataBody} bind:error />
	{:else if contentType === 'text/plain'}
		<PlainTextBodyInput bind:body={textBody} bind:error />
	{:else if contentType === 'application/xml'}
		<XMLBodyInput bind:body={textBody} bind:error />
	{:else if contentType === 'application/octet-stream'}
		<ByteBodyInput bind:body={textBody} bind:error />
	{:else}
		<p>Selected Content Type is not supported</p>
	{/if}
	{#if error}
		<div class="error">{error}</div>
	{/if}
</section>
