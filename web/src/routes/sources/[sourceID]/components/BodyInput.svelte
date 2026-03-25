<script lang="ts">
	import Select from '$lib/components/ui/Select.svelte';
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

<section class="flex flex-col gap-4">
	<div>
		<label for="content-type" class="text-sm font-medium text-fg">Content type</label>
		<Select
			id="content-type"
			bind:value={contentType}
			options={[
				{ value: 'application/json', label: 'JSON' },
				{ value: 'application/x-www-form-urlencoded', label: 'Form URL Encoded' },
				{ value: 'multipart/form-data', label: 'Multipart Form Data' },
				{ value: 'text/plain', label: 'Plain Text' },
				{ value: 'application/xml', label: 'XML' },
				{ value: 'application/octet-stream', label: 'Binary Data' }
			]}
			class="mt-2 w-full rounded-md border border-border bg-surface px-4 py-3 text-sm text-fg shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
		/>
	</div>

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
		<p class="text-sm text-muted">Selected content type is not supported.</p>
	{/if}
	{#if error}
		<div class="rounded-md border border-error bg-surface px-4 py-3 text-sm text-error">{error}</div>
	{/if}
</section>
