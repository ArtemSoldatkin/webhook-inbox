<script lang="ts">
	import PlainTextBodyInput from '../../routes/sources/[sourceID]/components/body/inputs/PlainTextBodyInput.svelte';

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

<PlainTextBodyInput bind:body bind:error />

<output data-testid="body-state">{body}</output>
<output data-testid="error-state">{error ?? ''}</output>
