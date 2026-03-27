<script lang="ts">
	import Alert from '$lib/components/ui/Alert.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Checkbox from '$lib/components/ui/Checkbox.svelte';
	import DateInput from '$lib/components/ui/DateInput.svelte';
	import Eyebrow from '$lib/components/ui/Eyebrow.svelte';
	import FileInput from '$lib/components/ui/FileInput.svelte';
	import NumberInput from '$lib/components/ui/NumberInput.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import TextInput from '$lib/components/ui/TextInput.svelte';
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
				case 'date':
					field.value = '';
					break;
				case 'number':
					field.value = null;
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
			case 'number':
				field.value = null;
				break;
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
			<Eyebrow as="label">
				Type
				<Select
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
			</Eyebrow>
		</div>

		<div>
			<Eyebrow as="label">
				Name
				<TextInput placeholder="Enter name" bind:value={field.name} />
			</Eyebrow>
		</div>

		<div>
			<Eyebrow>Value</Eyebrow>
			{#if field.name.trim() === ''}
				<Alert variant="warning" class="mt-1 shadow-none">Name is required.</Alert>
			{:else if field.type === 'text'}
				<TextInput
					name={field.name}
					placeholder="Enter text"
					bind:value={field.value}
					class="mt-1 w-full"
				/>
			{:else if field.type === 'number'}
				<NumberInput
					name={field.name}
					placeholder="Enter number"
					bind:value={field.value}
					class="mt-1 w-full"
				/>
			{:else if field.type === 'checkbox'}
				<label
					class="mt-1 flex items-center gap-3 rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm"
				>
					<Checkbox name={field.name} bind:value={field.value} />
					<span>Checked</span>
				</label>
			{:else if field.type === 'file'}
				<FileInput
					name={field.name}
					bind:value={field.value}
					class="mt-1 block w-full"
				/>
			{:else if field.type === 'date'}
				<DateInput
					name={field.name}
					bind:value={field.value}
					class="mt-1 w-full"
				/>
			{:else}
				<Alert variant="warning" class="mt-1 shadow-none">Unsupported input type.</Alert>
			{/if}
		</div>
	</div>

	<div class="mt-4 flex justify-end">
		<Button type="button" onclick={removeValue} variant="secondary">Clear field value</Button>
	</div>
</section>
