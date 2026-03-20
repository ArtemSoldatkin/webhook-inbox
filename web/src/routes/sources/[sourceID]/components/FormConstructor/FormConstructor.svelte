<script lang="ts">
	import type { FormField } from '../types';
	import InputConstructor from './InputConstructor.svelte';

	type Props = {
		fields: FormField[];
	};

	let { fields = $bindable() }: Props = $props();

	function handleAddField() {
		fields = [...fields, { type: 'text', name: '', value: '' }];
	}
</script>

{#if fields.length === 0}
	<p>No fields added yet. Click "Add field" to start building your form.</p>
{/if}
{#each fields as _, index (index)}
	<InputConstructor bind:field={fields[index]} />
{/each}
<button type="button" onclick={handleAddField}>Add field</button>
