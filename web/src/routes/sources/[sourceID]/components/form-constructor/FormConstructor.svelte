<script lang="ts">
	import type { FormField } from '../types';
	import InputConstructor from './InputConstructor.svelte';

	/** Props for the form field constructor. */
	type Props = {
		/** Bound list of form fields in the constructor. */
		fields: FormField[];
	};

	let { fields = $bindable() }: Props = $props();

	/** Adds a new empty field to the multipart form builder. */
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
