<script lang="ts">
	import FormDataBodyInput from '../../routes/sources/[sourceID]/components/body/inputs/FormDataBodyInput.svelte';

	let {
		initialBody = new FormData(),
		initialError = null
	}: {
		initialBody?: FormData;
		initialError?: string | null;
	} = $props();

	let body = $state(initialBody);
	let error = $state<string | null>(initialError);

	function stringifyFormData(formData: FormData): string {
		return JSON.stringify(
			Array.from(formData.entries()).map(([key, value]) => [
				key,
				value instanceof File ? `file:${value.name}` : value
			])
		);
	}
</script>

<FormDataBodyInput bind:body bind:error />

<output data-testid="formdata-state">{stringifyFormData(body)}</output>
<output data-testid="error-state">{error ?? ''}</output>
