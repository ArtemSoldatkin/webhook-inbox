<script lang="ts">
	import JSONBodyInput from '../../routes/sources/[sourceID]/components/body/inputs/JSONBodyInput.svelte';

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

<JSONBodyInput bind:body bind:error />

<output data-testid="body-state">{body}</output>
<output data-testid="error-state">{error ?? ''}</output>
