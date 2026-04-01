<script lang="ts">
	import XMLBodyInput from '../../routes/sources/[sourceID]/components/body/inputs/XMLBodyInput.svelte';

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

<XMLBodyInput bind:body bind:error />

<output data-testid="body-state">{body}</output>
<output data-testid="error-state">{error ?? ''}</output>
