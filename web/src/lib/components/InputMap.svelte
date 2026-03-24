<script lang="ts">
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

<ul>
	{#each Object.keys(map) as mapKey (mapKey)}
		<li>
			<p>{mapKey}</p>
			<input type="text" placeholder="Value" bind:value={map[mapKey]} {disabled} />
			<button type="button" {disabled} onclick={() => removeKey(mapKey)}>Remove</button>
		</li>
	{/each}
</ul>
<input type="text" placeholder="Key" bind:value={key} {disabled} />
<input type="text" placeholder="Value" bind:value {disabled} />
<button type="button" {disabled} onclick={() => addKeyValue()}>Add</button>
