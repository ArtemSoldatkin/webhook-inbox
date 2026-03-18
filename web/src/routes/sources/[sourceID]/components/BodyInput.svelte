<script lang="ts">
	import type { ContentType } from '$lib/types';
	import ByteBodyInput from './BodyTypeInputs/ByteBodyInput.svelte';
	import FormDataBodyInput from './BodyTypeInputs/FormDataBodyInput.svelte';
	import FormUrlEncodedBodyInput from './BodyTypeInputs/FormUrlEncodedBodyInput.svelte';
	import JSONBodyInput from './BodyTypeInputs/JSONBodyInput.svelte';
	import PlainTextBodyInput from './BodyTypeInputs/PlainTextBodyInput.svelte';
	import XMLBodyInput from './BodyTypeInputs/XMLBodyInput.svelte';

	type Props = {
		body: string;
		contentType: ContentType;
	};

	let { body, contentType }: Props = $props();

	let error = $state<string | null>(null);
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
		<JSONBodyInput bind:body bind:error />
	{:else if contentType === 'application/x-www-form-urlencoded'}
		<FormUrlEncodedBodyInput bind:body bind:error />
	{:else if contentType === 'multipart/form-data'}
		<FormDataBodyInput bind:body bind:error />
	{:else if contentType === 'text/plain'}
		<PlainTextBodyInput bind:body bind:error />
	{:else if contentType === 'application/xml'}
		<XMLBodyInput bind:body bind:error />
	{:else if contentType === 'application/octet-stream'}
		<ByteBodyInput bind:body bind:error />
	{:else}
		<p>Selected Content Type is not supported</p>
	{/if}
	{#if error}
		<div class="error">{error}</div>
	{/if}
</section>
