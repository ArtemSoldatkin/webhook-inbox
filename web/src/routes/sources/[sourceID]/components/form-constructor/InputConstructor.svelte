<script lang="ts">
	import type { FormField } from '../types';

	/** Props for a single dynamic form field editor. */
	type Props = {
		/** Bound field being edited by this constructor row. */
		field: FormField;
	};

	let { field = $bindable() }: Props = $props();

	/** Previous field type used to detect type changes. */
	let lastType = field.type;
	$effect(() => {
		if (field.type !== lastType) {
			switch (field.type) {
				case 'text':
				case 'number':
				case 'date':
					field.value = '';
					break;
				case 'checkbox':
					field.value = false;
					break;
				case 'file':
					field.value = null;
					break;
			}
			lastType = field.type;
		}
	});
</script>

<label>
	Type:
	<select bind:value={field.type}>
		<option value="text">Text</option>
		<option value="number">Number</option>
		<option value="checkbox">Checkbox</option>
		<option value="file">File</option>
		<option value="date">Date</option>
	</select>
</label>
<label>
	Name:
	<input type="text" placeholder="Enter name" bind:value={field.name} />
</label>
{#if field.name.trim() === ''}
	<p>Name is required</p>
{:else}
	<label>
		Value:
		{#if field.type === 'text'}
			<input type="text" name={field.name} placeholder="Enter text" bind:value={field.value} />
		{:else if field.type === 'number'}
			<input type="number" name={field.name} placeholder="Enter number" bind:value={field.value} />
		{:else if field.type === 'checkbox'}
			<input
				type="checkbox"
				name={field.name}
				placeholder="Enter boolean"
				bind:checked={field.value}
			/>
		{:else if field.type === 'file'}
			<input type="file" name={field.name} placeholder="Enter file" bind:files={field.value} />
		{:else if field.type === 'date'}
			<input type="date" name={field.name} placeholder="Enter date" bind:value={field.value} />
		{:else}
			<p>Unsupported input type</p>
		{/if}
	</label>
{/if}
