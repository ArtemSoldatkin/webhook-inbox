<script lang="ts">
	import InputConstructor from '../../routes/sources/[sourceID]/components/body/inputs/form-constructor/InputConstructor.svelte';
	import type { FormField } from '../../routes/sources/[sourceID]/components/body/inputs/types';

	let { initialField }: { initialField: FormField } = $props();

	let field = $state<FormField>(initialField);
</script>

<InputConstructor bind:field />

<output data-testid="field-state">{JSON.stringify(field)}</output>
