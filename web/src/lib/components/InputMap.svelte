<script lang="ts">
	type Props = {
		map: Record<string, string>;
		disabled?: boolean;
	};

	let { map = $bindable(), disabled }: Props = $props();

	let key = $state('');
	let value = $state('');

	function addKeyValue() {
		if (key.trim() === '' || value.trim() === '' || map[key]) {
			return;
		}
		map = { ...map, [key]: value };
		key = '';
		value = '';
	}

	function removeKey(keyToRemove: string) {
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
