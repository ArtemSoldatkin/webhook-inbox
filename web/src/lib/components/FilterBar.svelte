<script lang="ts">
	export let searchQuery: string = '';
	export let onSearch: () => void;

	export let filterOptions: string[] = [];
	export let filter: string = '*';
</script>

<section>
	<div>
		<input type="text" placeholder="Search..." bind:value={searchQuery} />
		<button type="button" on:click={onSearch}>Search</button>
	</div>
	{#if filterOptions.length > 0}
		<div>
			<select bind:value={filter}>
				<option value="*">All</option>
				{#each filterOptions as category (category)}
					<option value={category}>{category}</option>
				{/each}
			</select>
		</div>
	{/if}
</section>
