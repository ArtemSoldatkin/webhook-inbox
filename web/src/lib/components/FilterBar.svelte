<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Eyebrow from '$lib/components/ui/Eyebrow.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import Icon from '@iconify/svelte';
	import Input from './ui/Input.svelte';

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
			<Eyebrow as="label">
				Search
				<Input type="text" class="w-full mt-1" placeholder="Search..." bind:value={searchInput} />
			</Eyebrow>
		</div>
		<Button type="button" class="py-3 border border-transparent" onclick={handleSearch}
			>Search</Button
		>
	</div>

	<div class="flex flex-col gap-4 sm:flex-row sm:items-end">
		{#if filterOptions}
			<div class="min-w-0 sm:min-w-44">
				<Eyebrow as="label">
					Filter by {filterName ?? 'category'}
					<Select
						bind:value={filter}
						options={[
							{ value: '*', label: 'All' },
							...(filterOptions?.map((category) => ({ value: category, label: category })) ?? [])
						]}
						class="w-full mt-1"
					/>
				</Eyebrow>
			</div>
		{/if}

		<div>
			<Eyebrow>Sort</Eyebrow>
			<Button
				type="button"
				class="py-3 mt-1 border border-transparent"
				onclick={toggleSortDirection}
				aria-label="Toggle sort direction"
				variant="secondary"
			>
				{#if sortDirection === 'ASC'}
					<Icon icon="mdi:sort-ascending" class="text-xl" />
				{:else}
					<Icon icon="mdi:sort-descending" class="text-xl" />
				{/if}
			</Button>
		</div>
	</div>
</section>
