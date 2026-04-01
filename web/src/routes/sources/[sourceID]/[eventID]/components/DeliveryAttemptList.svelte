<script lang="ts">
	import { fetchPaginatedData } from '$lib/api';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import PageSizeSelector from '$lib/components/PageSizeSelector.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Eyebrow from '$lib/components/ui/Eyebrow.svelte';
	import KeyValueList from '$lib/components/ui/KeyValueList.svelte';
	import SectionHeader from '$lib/components/ui/SectionHeader.svelte';
	import { parseDeliveryAttemptDTO } from '$lib/dto-parsers';
	import { type DeliveryAttemptDTO } from '$lib/types';
	import { untrack } from 'svelte';

	type Props = {
		/** Source id that owns the current event. */
		sourceID: string;

		/** Event id whose delivery attempts are shown. */
		eventID: string;
	};

	/** Filters applied to the delivery attempts list. */
	type DeliveryAttemptFilters = {
		/** Source id that owns the current event. */
		sourceID: string;

		/** Event id whose delivery attempts are shown. */
		eventID: string;

		/** Requested page size for delivery attempt results. */
		pageSize: number;

		/** Free-text query applied to the delivery attempts endpoint. */
		searchQuery: string;

		/** Selected delivery state filter. Supported values are 'pending', 'in_flight', 'succeeded', 'failed', 'aborted', and '*' for no filtering. */
		filterState: string;

		/** Sort order used for delivery attempt results. */
		sortDirection: 'ASC' | 'DESC';
	};

	let { sourceID, eventID }: Props = $props();

	/** Loaded delivery attempts for the current event. */
	let data = $state<DeliveryAttemptDTO[]>([]);

	/** Tracks whether delivery attempts are loading. */
	let loading = $state(false);

	/** Holds the latest delivery attempt loading error. */
	let error = $state<string | null>(null);

	/** Requested page size for delivery attempts. */
	let pageSize = $state(20);

	/** Cursor used to fetch the next page of attempts. */
	let nextCursor = $state<string | null>(null);

	/** Indicates whether more delivery attempts are available. */
	let hasNext = $state(false);

	/** Free-text query applied to delivery attempts. */
	let searchQuery = $state('');

	/** Selected delivery state filter. */
	let filterState = $state('*');

	/** Available delivery attempt states for filtering. */
	const filterStateOptions = ['pending', 'in_flight', 'succeeded', 'failed', 'aborted'];

	/** Sort order used for delivery attempts. */
	let sortDirection = $state<'ASC' | 'DESC'>('DESC');

	/**
	 * Collects the current filters into a single object for easier passing to fetch functions.
	 *
	 * @returns Current delivery attempt filters.
	 */
	function getCurrentFilters(): DeliveryAttemptFilters {
		return {
			sourceID,
			eventID,
			pageSize,
			searchQuery,
			filterState,
			sortDirection
		};
	}

	/** Loads the next page of delivery attempts.
	 *
	 * @param filters - Filters to apply when fetching delivery attempts.
	 */
	async function fetchDeliveryAttempts(filters: DeliveryAttemptFilters): Promise<void> {
		loading = true;
		error = null;
		try {
			const urlSearchParams: Record<string, string> = {};
			if (filters.searchQuery) urlSearchParams.search = filters.searchQuery;
			if (filters.filterState) urlSearchParams.filter_state = filters.filterState;
			if (filters.sortDirection) urlSearchParams.sort_direction = filters.sortDirection;

			const result = await fetchPaginatedData(
				`/api/sources/${filters.sourceID}/events/${filters.eventID}/delivery-attempts`,
				filters.pageSize,
				nextCursor,
				urlSearchParams
			);
			data = [...data, ...result.data.map(parseDeliveryAttemptDTO)];
			nextCursor = result.next_cursor;
			hasNext = result.has_next;
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching event details:', err);
		} finally {
			loading = false;
		}
	}

	/** Resets delivery attempt pagination and fetches from the start.
	 *
	 * @param filters - Filters to apply when fetching delivery attempts.
	 */
	async function resetAndFetchDeliveryAttempts(filters: DeliveryAttemptFilters): Promise<void> {
		data = [];
		nextCursor = null;
		hasNext = false;
		await fetchDeliveryAttempts(filters);
	}

	/** Refreshes the delivery attempt list with the current filters when the user clicks the "Refresh" button. */
	function handleRefresh(): void {
		void resetAndFetchDeliveryAttempts(getCurrentFilters());
	}

	/** Loads the next page of delivery attempts when the user clicks the "Load More" button. */
	function handleLoadMore(): void {
		void fetchDeliveryAttempts(getCurrentFilters());
	}

	$effect(() => {
		const filters = getCurrentFilters();

		untrack(() => {
			void resetAndFetchDeliveryAttempts(filters);
		});
	});
</script>

<section class="rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8">
	<div class="flex flex-col gap-6">
		<SectionHeader
			eyebrow="Delivery attempts"
			title="Delivery history for this event"
			description="Review retries, status codes, and error information for each outbound delivery attempt."
			titleAs="h3"
		>
			{#snippet actions()}
				<Button onclick={handleRefresh} disabled={loading} variant="secondary">
					Refresh Delivery Attempts
				</Button>
			{/snippet}
		</SectionHeader>

		<FilterBar
			bind:searchQuery
			bind:filter={filterState}
			bind:sortDirection
			filterName="state"
			filterOptions={filterStateOptions}
		/>

		{#if loading}
			<Alert>Loading delivery attempts...</Alert>
		{:else if error}
			<Alert variant="error" title="Error" class="bg-surface">{error}</Alert>
		{:else if data.length === 0}
			<Alert>No delivery attempts found for this event.</Alert>
		{:else}
			<ul class="grid gap-4">
				{#each data as attempt (attempt.id)}
					<li>
						<article class="rounded-lg border border-border bg-surface p-5 shadow-sm">
							<div class="flex flex-col gap-5">
								<div class="flex flex-wrap items-center gap-3">
									<h4 class="text-xl font-semibold tracking-tight text-fg">
										Attempt ID: {attempt.id}
									</h4>
									<Badge variant="neutral" appearance="soft">{attempt.state}</Badge>
								</div>

								<div
									class="grid gap-4 border-t border-border-muted pt-4 sm:grid-cols-2 xl:grid-cols-3"
								>
									<div>
										<Eyebrow>Event ID</Eyebrow>
										<p class="mt-2 break-all text-sm text-fg">{attempt.event_id}</p>
									</div>
									<div>
										<Eyebrow>Attempt number</Eyebrow>
										<p class="mt-2 text-sm text-fg">{attempt.attempt_number}</p>
									</div>
									<div>
										<Eyebrow>Status code</Eyebrow>
										<p class="mt-2 text-sm text-fg">{attempt.status_code ?? 'N/A'}</p>
									</div>
								</div>

								<div class="grid gap-4 border-t border-border-muted pt-4 lg:grid-cols-2">
									<div>
										<Eyebrow>Errors</Eyebrow>
										<KeyValueList
											items={[
												{ label: 'Error type', value: attempt.error_type },
												{ label: 'Error message', value: attempt.error_message }
											]}
										/>
									</div>
									<div>
										<Eyebrow>Timing</Eyebrow>
										<KeyValueList
											items={[
												{
													label: 'Started at',
													value: attempt.started_at
														? new Date(attempt.started_at).toLocaleString()
														: 'N/A'
												},
												{
													label: 'Finished at',
													value: attempt.finished_at
														? new Date(attempt.finished_at).toLocaleString()
														: 'N/A'
												},
												{
													label: 'Created at',
													value: new Date(attempt.created_at).toLocaleString()
												},
												{
													label: 'Next attempt at',
													value: attempt.next_attempt_at
														? new Date(attempt.next_attempt_at).toLocaleString()
														: 'N/A'
												}
											]}
										/>
									</div>
								</div>
							</div>
						</article>
					</li>
				{/each}
			</ul>

			{#if hasNext}
				<div class="flex justify-center pt-2">
					<Button onclick={handleLoadMore} disabled={loading} variant="secondary">Load More</Button>
				</div>
			{/if}
		{/if}

		<div class="border-t border-border-muted pt-4">
			<PageSizeSelector bind:pageSize />
		</div>
	</div>
</section>
