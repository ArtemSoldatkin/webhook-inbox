<script lang="ts">
	type Props = {
		/** Bound page size for paginated lists. */
		pageSize: number;
	};

	let { pageSize = $bindable() }: Props = $props();
</script>

<div class="flex items-center gap-3">
	<label for="page-size" class="text-sm font-medium text-muted">Page size</label>
	<select
		id="page-size"
		bind:value={pageSize}
		class="rounded-md border border-border bg-surface px-4 py-2 text-sm text-fg shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
	>
		<option value={1}>1</option>
		<option value={5}>5</option>
		<option value={10}>10</option>
		<option value={20}>20</option>
		<option value={50}>50</option>
		<option value={100}>100</option>
	</select>
</div>
