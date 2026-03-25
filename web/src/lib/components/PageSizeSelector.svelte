<script lang="ts">
	import Select from '$lib/components/ui/Select.svelte';

	type Props = {
		/** Bound page size for paginated lists. */
		pageSize: number;
	};

	let { pageSize = $bindable() }: Props = $props();
</script>

<div class="flex items-center gap-3">
	<label for="page-size" class="text-sm font-medium text-muted">Page size</label>
	<Select
		id="page-size"
		bind:value={pageSize}
		options={[
			{ value: 1, label: '1' },
			{ value: 5, label: '5' },
			{ value: 10, label: '10' },
			{ value: 20, label: '20' },
			{ value: 50, label: '50' },
			{ value: 100, label: '100' }
		]}
		class="rounded-md border border-border bg-surface px-4 py-2 text-sm text-fg shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
	/>
</div>
