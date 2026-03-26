<script lang="ts">
	import Alert from '$lib/components/ui/Alert.svelte';
	import Button from '$lib/components/ui/Button.svelte';

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
				<li
					class="rounded-md border border-border-muted bg-surface p-4 shadow-sm"
				>
					<div class="flex flex-col gap-3 sm:flex-row sm:items-start">
						<div class="min-w-0 flex-1">
							<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">Header key</p>
							<p class="mt-1 break-all text-sm font-medium text-fg">{mapKey}</p>
						</div>
						<div class="min-w-0 flex-1">
							<label class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
								Header value
								<input
									type="text"
									placeholder="Value"
									bind:value={map[mapKey]}
									{disabled}
									class="mt-1 w-full rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
								/>
							</label>
						</div>
						<div class="sm:pt-6">
							<Button type="button" {disabled} onclick={() => removeKey(mapKey)} variant="secondary">
								Remove
							</Button>
						</div>
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
				<label class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
					Key
					<input
						type="text"
						placeholder="Authorization"
						bind:value={key}
						{disabled}
						class="mt-1 w-full rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
					/>
				</label>
			</div>
			<div>
				<label class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
					Value
					<input
						type="text"
						placeholder="Bearer secret"
						bind:value
						{disabled}
						class="mt-1 w-full rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
					/>
				</label>
			</div>
			<Button type="button" {disabled} onclick={() => addKeyValue()}>
				Add
			</Button>
		</div>
	</div>
</div>
