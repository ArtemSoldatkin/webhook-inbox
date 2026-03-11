<script lang="ts">
	export let searchQuery: string = '';
	export let onSearch: () => void;
</script>

<section>
	<div>
		<input type="text" placeholder="Search..." bind:value={searchQuery} />
		<button type="button" on:click={onSearch}>Search</button>
	</div>
</section>
