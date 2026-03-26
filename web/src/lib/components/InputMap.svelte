<script lang="ts">
	import Alert from '$lib/components/ui/Alert.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Eyebrow from '$lib/components/ui/Eyebrow.svelte';
	import Input from './ui/Input.svelte';

	type Props = {
		/** Bound key-value pairs being edited. */
		map: Record<string, string>;

		/** Disables all map inputs when true. */
		disabled?: boolean;
	};

	let { map = $bindable(), disabled }: Props = $props();

	/** Draft key for the next map entry. */
	let key = $state('');

	/** Draft value for the next map entry. */
	let value = $state('');

	/** Adds the current draft entry when both fields are valid. */
	function addKeyValue(): void {
		if (key.trim() === '' || value.trim() === '' || map[key]) {
			return;
		}
		map = { ...map, [key]: value };
		key = '';
		value = '';
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
</script>

<div class="flex flex-col gap-4">
	{#if Object.keys(map).length > 0}
		<ul class="flex flex-col gap-3">
			{#each Object.keys(map) as mapKey (mapKey)}
				<li class="rounded-md border border-border-muted bg-surface p-4">
					<div class="grid gap-3 sm:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_auto] sm:items-end">
						<div>
							<Eyebrow>Header key</Eyebrow>
							<p class="mt-1 py-3 break-all text-sm font-medium text-fg">{mapKey}</p>
						</div>
						<div>
							<Eyebrow as="label"
								>Header value
								<Input type="text" placeholder="Value" bind:value={map[mapKey]} {disabled} />
							</Eyebrow>
						</div>
						<Button
							type="button"
							variant="secondary"
							class="py-3 border border-transparent"
							{disabled}
							onclick={() => removeKey(mapKey)}
						>
							x
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
			<div>
				<Eyebrow as="label"
					>Key
					<Input type="text" placeholder="Key" bind:value={key} {disabled} />
				</Eyebrow>
			</div>
			<div>
				<Eyebrow as="label"
					>Value
					<Input type="text" placeholder="Value" bind:value {disabled} />
				</Eyebrow>
			</div>
			<Button
				type="button"
				class="py-3 border border-transparent"
				disabled={disabled || !key || key in map}
				onclick={() => addKeyValue()}
			>
				+
			</Button>
		</div>
	</div>
</div>
