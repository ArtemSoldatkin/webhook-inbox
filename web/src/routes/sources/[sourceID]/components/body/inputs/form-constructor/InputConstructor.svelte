<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import type { FormField } from '../types';

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

	function removeValue(): void {
		switch (field.type) {
			case 'checkbox':
				field.value = false;
				break;
			case 'file':
				field.value = null;
				break;
			default:
				field.value = '';
		}
	}
</script>

<section class="rounded-md border border-border-muted bg-surface p-4">
	<div class="grid gap-4 lg:grid-cols-[minmax(0,0.24fr)_minmax(0,0.3fr)_minmax(0,1fr)]">
		<div>
			<label for={`field-type-${field.name || 'new'}`} class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
				Type
			</label>
			<Select
				id={`field-type-${field.name || 'new'}`}
				bind:value={field.type}
				options={[
					{ value: 'text', label: 'Text' },
					{ value: 'number', label: 'Number' },
					{ value: 'checkbox', label: 'Checkbox' },
					{ value: 'file', label: 'File' },
					{ value: 'date', label: 'Date' }
				]}
				class="mt-1 w-full rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
			/>
		</div>

		<div>
			<label for={`field-name-${field.name || 'new'}`} class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
				Name
			</label>
			<input
				id={`field-name-${field.name || 'new'}`}
				type="text"
				placeholder="Enter name"
				bind:value={field.name}
				class="mt-1 w-full rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
			/>
		</div>

		<div>
			<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">Value</p>
			{#if field.name.trim() === ''}
				<div class="mt-1 rounded-md border border-warning bg-elevated px-4 py-3 text-sm text-warning">
					Name is required.
				</div>
			{:else if field.type === 'text'}
				<input
					type="text"
					name={field.name}
					placeholder="Enter text"
					bind:value={field.value}
					class="mt-1 w-full rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
				/>
			{:else if field.type === 'number'}
				<input
					type="number"
					name={field.name}
					placeholder="Enter number"
					bind:value={field.value}
					class="mt-1 w-full rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
				/>
			{:else if field.type === 'checkbox'}
				<label class="mt-1 flex items-center gap-3 rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm">
					<input type="checkbox" name={field.name} bind:checked={field.value} />
					<span>Checked</span>
				</label>
			{:else if field.type === 'file'}
				<input
					type="file"
					name={field.name}
					bind:files={field.value}
					class="mt-1 block w-full rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm"
				/>
			{:else if field.type === 'date'}
				<input
					type="date"
					name={field.name}
					bind:value={field.value}
					class="mt-1 w-full rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
				/>
			{:else}
				<div class="mt-1 rounded-md border border-warning bg-elevated px-4 py-3 text-sm text-warning">
					Unsupported input type.
				</div>
			{/if}
		</div>
	</div>

	<div class="mt-4 flex justify-end">
		<Button type="button" onclick={removeValue} variant="secondary">Clear field value</Button>
	</div>
</section>
