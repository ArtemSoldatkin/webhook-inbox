<script lang="ts">
	export let searchQuery: string = '';
	export let onSearch: (query: string) => void;

	function handleSearch() {
		onSearch(searchQuery);
	}
</script>

<section>
	<div>
		<input type="text" placeholder="Search..." bind:value={searchQuery} />
		<button type="button" on:click={handleSearch}>Search</button>
	</div>
</section>
