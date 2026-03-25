<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Select from '$lib/components/ui/Select.svelte';

	type Props = {
		/** Bound search text for the list. */
		searchQuery: string;

		/** Active sort direction for the list. */
		sortDirection: 'ASC' | 'DESC';

		/** Label shown for the optional filter select. */
		filterName?: string;

		/** Available values for the optional filter select. */
		filterOptions?: string[];

		/** Currently selected filter value. */
		filter?: string;
	};

	let {
		searchQuery = $bindable(),
		sortDirection = $bindable(),
		filterName,
		filterOptions,
		filter = $bindable()
	}: Props = $props();

	let searchInput = $state('');

	/** Applies the current search input as a filter. */
	function handleSearch(): void {
		searchQuery = searchInput;
	}

	/** Flips the active sort direction between ascending and descending. */
	function toggleSortDirection(): void {
		sortDirection = sortDirection === 'ASC' ? 'DESC' : 'ASC';
	}
</script>

<section class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
	<div class="flex flex-1 flex-col gap-4 sm:flex-row sm:items-end">
		<div class="min-w-0 flex-1">
			<label for="search-query" class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
				Search
			</label>
			<input
				id="search-query"
				type="text"
				placeholder="Search..."
				bind:value={searchInput}
				class="mt-1 w-full rounded-md border border-border bg-surface px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
			/>
		</div>
		<Button type="button" onclick={handleSearch}>Search</Button>
	</div>

	<div class="flex flex-col gap-4 sm:flex-row sm:items-end">
		{#if filterOptions}
			<div class="min-w-0 sm:min-w-44">
				<label for="filter" class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
					Filter by {filterName ?? 'category'}
				</label>
				<Select
					id="filter"
					bind:value={filter}
					options={[
						{ value: '*', label: 'All' },
						...(filterOptions?.map((category) => ({ value: category, label: category })) ?? [])
					]}
					class="mt-1 w-full rounded-md border border-border bg-surface px-4 py-3 text-sm text-fg shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
				/>
			</div>
		{/if}

		<div>
			<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">Sort</p>
			<Button type="button" onclick={toggleSortDirection} aria-label="Toggle sort direction" variant="secondary">
				{sortDirection === 'ASC' ? 'Ascending' : 'Descending'}
			</Button>
		</div>
	</div>
</section>
