<script lang="ts">
	export let pageSize: number = 20;
</script>

<select bind:value={pageSize}>
	<option value={1}>1</option>
	<option value={5}>5</option>
	<option value={10}>10</option>
	<option value={20}>20</option>
	<option value={50}>50</option>
	<option value={100}>100</option>
</select>
