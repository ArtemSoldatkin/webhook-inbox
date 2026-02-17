<script lang="ts">
	let key: string = '';
	let value: string = '';

	export let json: Record<string, string> = {};
	export let disabled: boolean = false;

	function addKeyValue() {
		if (key.trim() === '' || value.trim() === '' || json[key]) {
			return;
		}
		json = { ...json, [key]: value };
		key = '';
		value = '';
	}

	function removeKey(keyToRemove: string) {
		delete json[keyToRemove];
		json = { ...json };
	}
</script>

<ul>
	{#each Object.keys(json) as jsonKey}
		<li>
			<p>{jsonKey}</p>
			<input type="text" placeholder="Value" bind:value={json[jsonKey]} {disabled} />
			<button type="button" {disabled} on:click={() => removeKey(jsonKey)}>Remove</button>
		</li>
	{/each}
</ul>
<input type="text" placeholder="Key" bind:value={key} {disabled} />
<input type="text" placeholder="Value" bind:value {disabled} />
<button type="button" {disabled} on:click={() => addKeyValue()}>Add</button>
