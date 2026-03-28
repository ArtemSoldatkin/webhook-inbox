<script lang="ts">
	import Alert from '$lib/components/ui/Alert.svelte';
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
	<Alert class="bg-surface shadow-none">
		No fields added yet. Click "Add field" to start building your form.
	</Alert>
{/if}
<div class="flex flex-col gap-4">
	<!-- eslint-disable-next-line @typescript-eslint/no-unused-vars -->
	{#each fields as _, index (index)}
		<InputConstructor bind:field={fields[index]} />
	{/each}
	<div class="flex justify-end">
		<Button type="button" onclick={handleAddField}>Add field</Button>
	</div>
</div>
