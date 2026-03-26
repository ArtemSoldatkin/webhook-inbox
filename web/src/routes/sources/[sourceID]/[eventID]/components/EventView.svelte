<script lang="ts">
	import Badge from '$lib/components/ui/Badge.svelte';
	import DisplayMapOfStringArrays from '$lib/components/DisplayMapOfStringArrays.svelte';
	import { parseEventDTO } from '$lib/dto-parsers';
	import { type EventDTO } from '$lib/types';
	import { untrack } from 'svelte';
	import BodyView from '../../components/BodyView.svelte';
	import DeliveryAttemptList from './DeliveryAttemptList.svelte';

	type Props = {
		/** Source id owning the current event. */
		sourceID: string;

		/** Event id shown on the details page. */
		eventID: string;
	};

	/** Filters applied to the event details view. */
	type EventFilters = {
		/** Source id owning the current event. */
		sourceID: string;
		/** Event id shown on the details page. */
		eventID: string;
	};

	let { sourceID, eventID }: Props = $props();

	/** Loaded event details. */
	let data = $state<EventDTO | null>(null);

	/** Tracks whether the event request is in flight. */
	let loading = $state(false);

	/** Holds the latest event loading error. */
	let error = $state<string | null>(null);

	/** Collects the current filters into a single object for easier passing to fetch functions. */
	function getCurrentFilters(): EventFilters {
		return {
			sourceID,
			eventID
		};
	}

	/** Fetches the current event and normalizes its payload.
	 *
	 * @param filters - Filters to apply when fetching the event.
	 */
	async function fetchEventDetails(filters: EventFilters): Promise<void> {
		loading = true;
		error = null;
		try {
			const response = await fetch(`/api/sources/${filters.sourceID}/events/${filters.eventID}`);
			if (!response.ok) {
				throw new Error(`Failed to fetch event details: ${response.statusText}`);
			}
			const rawData = await response.json();
			data = parseEventDTO(rawData);
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching event details:', err);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		const filters = getCurrentFilters();

		untrack(() => {
			void fetchEventDetails(filters);
		});
	});
</script>

{#if loading}
	<div class="rounded-md border border-border-muted bg-elevated px-4 py-6 text-sm text-muted">
		Loading event details...
	</div>
{:else if error}
	<div class="rounded-md border border-error bg-surface px-4 py-3 text-sm text-error">
		Error: {error}
	</div>
{:else if data}
	<div class="flex flex-col gap-8">
		<section class="rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8">
			<div class="flex flex-col gap-6">
				<div class="flex flex-col gap-3">
					<div class="flex flex-wrap items-center gap-3">
						<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">Captured event</p>
						<Badge variant="neutral" appearance="soft" class="bg-elevated">{data.method}</Badge>
					</div>
					<h2 class="text-3xl font-semibold tracking-tight text-fg">Event ID: {data.id}</h2>
				</div>

				<div class="grid gap-4 sm:grid-cols-2">
					<div class="rounded-md border border-border-muted bg-elevated p-4">
						<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">Source ID</p>
						<p class="mt-2 break-all text-sm leading-6 text-fg">{data.source_id}</p>
					</div>
					<div class="rounded-md border border-border-muted bg-elevated p-4">
						<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
							Deduplication hash
						</p>
						<p class="mt-2 break-all text-sm leading-6 text-fg">{data.dedup_hash ?? 'N/A'}</p>
					</div>
				</div>

				<div class="grid gap-4 lg:grid-cols-2">
					<DisplayMapOfStringArrays title="Query Parameters" data={data.query_params ?? {}} />
					<DisplayMapOfStringArrays title="Raw Headers" data={data.raw_headers ?? {}} />
				</div>

				<BodyView body={data.body} contentType={data.body_content_type} />
			</div>
		</section>

		<DeliveryAttemptList {sourceID} {eventID} />
	</div>
{:else}
	<div class="rounded-md border border-border-muted bg-elevated px-4 py-6 text-sm text-muted">
		No details found for this event.
	</div>
{/if}
