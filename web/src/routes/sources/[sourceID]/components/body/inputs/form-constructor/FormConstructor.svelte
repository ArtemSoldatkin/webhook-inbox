<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import type { FormField } from '../types';
	import InputConstructor from './InputConstructor.svelte';

	type Props = {
		/** Bound list of form fields in the constructor. */
		fields: FormField[];
	};

	let { fields = $bindable() }: Props = $props();

	/** Adds a new empty field to the multipart form builder. */
	function handleAddField(): void {
		fields = [...fields, { type: 'text', name: '', value: '' }];
	}
</script>

{#if fields.length === 0}
	<p>No fields added yet. Click "Add field" to start building your form.</p>
{/if}
<!-- eslint-disable-next-line @typescript-eslint/no-unused-vars -->
{#each fields as _, index (index)}
	<InputConstructor bind:field={fields[index]} />
{/each}
<Button type="button" onclick={handleAddField}>Add field</Button>
