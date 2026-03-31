<script lang="ts">
	import Alert from '$lib/components/ui/Alert.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Eyebrow from '$lib/components/ui/Eyebrow.svelte';
	import { cx } from '$lib/utils/cx';
	import Icon from '@iconify/svelte';
	import TextInput from './ui/TextInput.svelte';

	type Props = {
		/** Bound key-value pairs being edited. */
		map: Record<string, string>;

		/** Additional CSS classes to apply to the component's root element. */
		class?: string;

		/** Disables all map inputs when true. */
		disabled?: boolean;
	};

	let { map = $bindable(), class: className, disabled = false }: Props = $props();

	/** Draft key for the next map entry. */
	let key = $state('');

	/** Draft value for the next map entry. */
	let value = $state('');

	/** Error message for the current draft entry, shown when the key is empty or already exists in the map. */
	let inputError = $state<string | null>(null);

	/** Adds the current draft entry when both fields are valid. */
	function addKeyValue(): void {
		if (key.trim() === '' || value.trim() === '' || map[key]) {
			return;
		}
		map = { ...map, [key]: value };
		key = '';
		value = '';
		inputError = null;
	}

	/**
	 * Removes a single entry from the bound map.
	 *
	 * @param keyToRemove - Map key to delete.
	 */
	function removeKey(keyToRemove: string): void {
		delete map[keyToRemove];
		map = { ...map };
	}

	/** Validates the draft key on blur, setting an error message if it's empty or already exists in the map. */
	function handleKeyBlur(): void {
		if (key.trim() === '') {
			inputError = 'Key cannot be empty.';
		} else if (key in map) {
			inputError = 'Key already exists in the map.';
		} else {
			inputError = null;
		}
	}
</script>

<div class={cx('flex flex-col gap-4', className)}>
	{#if Object.keys(map).length > 0}
		<ul class="flex flex-col gap-3" aria-label="Mapped key value pairs">
			{#each Object.keys(map) as mapKey (mapKey)}
				<li class="rounded-md border border-border-muted bg-surface p-4">
					<div class="grid gap-3 sm:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_auto] sm:items-end">
						<Eyebrow
							>Header key
							<p class="mt-1 py-3 break-all text-sm font-medium text-fg">{mapKey}</p>
						</Eyebrow>
						<Eyebrow as="label"
							>Header value
							<TextInput
								class="w-full mt-1"
								placeholder="Value"
								bind:value={map[mapKey]}
								{disabled}
							/>
						</Eyebrow>
						<Button
							type="button"
							variant="secondary"
							class="py-3 border border-transparent text-error"
							{disabled}
							onclick={() => removeKey(mapKey)}
							aria-label={`Remove ${mapKey}`}
						>
							<Icon icon="material-symbols:delete-rounded" class="text-xl" aria-hidden="true" />
						</Button>
					</div>
				</li>
			{/each}
		</ul>
	{:else}
		<Alert class="bg-surface shadow-none">No static headers added yet.</Alert>
	{/if}

	<div class="rounded-md border border-border-muted bg-surface p-4">
		<div class="grid gap-3 sm:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_auto] sm:items-end">
			<Eyebrow as="label"
				>Key
				<TextInput
					class="w-full mt-1"
					placeholder="Key"
					bind:value={key}
					{disabled}
					onblur={handleKeyBlur}
					aria-label="Key"
					aria-invalid={inputError ? 'true' : undefined}
				/>
			</Eyebrow>

			<Eyebrow as="label"
				>Value
				<TextInput
					class="w-full mt-1"
					placeholder="Value"
					bind:value
					{disabled}
					aria-label="Value"
				/>
			</Eyebrow>

			<Button
				type="button"
				class="py-3 border border-transparent"
				disabled={disabled || !key || key in map || !value}
				onclick={() => addKeyValue()}
				aria-label="Add key value pair"
			>
				<Icon icon="material-symbols:add-rounded" class="text-xl" aria-hidden="true" />
			</Button>
		</div>
		{#if inputError}
			<Alert variant="error" class="mt-3">
				{inputError}
			</Alert>
		{/if}
	</div>
</div>
