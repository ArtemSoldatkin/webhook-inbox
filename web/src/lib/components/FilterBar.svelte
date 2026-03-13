<script lang="ts">
	let {
		searchQuery = $bindable(),
		onSearch,
		filterName,
		filterOptions,
		filter = $bindable(),
		sortDirection = $bindable()
	} = $props<{
		searchQuery: string;
		onSearch: () => void;
		filterName?: string;
		filterOptions?: string[];
		filter?: string;
		sortDirection: 'ASC' | 'DESC';
	}>();

	function toggleSortDirection() {
		sortDirection = sortDirection === 'ASC' ? 'DESC' : 'ASC';
	}
</script>

<section>
	<div>
		<input type="text" placeholder="Search..." bind:value={searchQuery} />
		<button type="button" on:click={onSearch}>Search</button>
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
		<button type="button" on:click={toggleSortDirection} aria-label="Toggle sort direction">
			{sortDirection === 'ASC' ? '↑' : '↓'}
		</button>
	</div>
</section>
