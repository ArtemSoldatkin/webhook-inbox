<script lang="ts">
	import FormUrlEncodedBodyInput from '../../routes/sources/[sourceID]/components/body/inputs/FormUrlEncodedBodyInput.svelte';

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

<FormUrlEncodedBodyInput bind:body bind:error />

<output data-testid="body-state">{body}</output>
<output data-testid="error-state">{error ?? ''}</output>
