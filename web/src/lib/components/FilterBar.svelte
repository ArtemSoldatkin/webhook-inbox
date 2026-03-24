<script lang="ts">
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

<section>
	<div>
		<input type="text" placeholder="Search..." bind:value={searchInput} />
		<button type="button" onclick={handleSearch}>Search</button>
	</div>
	{#if filter && filterOptions}
		<div>
			<label for="filter"
				>Filter by {filterName ?? 'category'}:
				<select bind:value={filter}>
					<option value="*">All</option>
					{#each filterOptions as category (category)}
						<option value={category}>{category}</option>
					{/each}
				</select>
			</label>
		</div>
	{/if}
	<div>
		<button type="button" onclick={toggleSortDirection} aria-label="Toggle sort direction">
			{sortDirection === 'ASC' ? '↑' : '↓'}
		</button>
	</div>
</section>
