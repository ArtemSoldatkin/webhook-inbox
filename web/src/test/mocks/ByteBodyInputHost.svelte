<script lang="ts">
	import ByteBodyInput from '../../routes/sources/[sourceID]/components/body/inputs/ByteBodyInput.svelte';

	let {
		initialBody = '',
		initialError = null
	}: {
		initialBody?: string;
		initialError?: string | null;
	} = $props();

	let body = $state(initialBody);
	let error = $state<string | null>(initialError);
</script>

<ByteBodyInput bind:body bind:error />

<output data-testid="body-state">{body}</output>
<output data-testid="error-state">{error ?? ''}</output>
