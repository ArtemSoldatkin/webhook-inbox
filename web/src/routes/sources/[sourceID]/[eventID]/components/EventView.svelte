<script lang="ts">
	import Alert from '$lib/components/ui/Alert.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import Eyebrow from '$lib/components/ui/Eyebrow.svelte';
	import { stringArrayRecordToKeyValueItems } from '$lib/components/ui/key-value-list';
	import KeyValueList from '$lib/components/ui/KeyValueList.svelte';
	import SectionHeader from '$lib/components/ui/SectionHeader.svelte';
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
				let errorMsg = `Request failed: ${response.statusText}`;
				try {
					const errorJson = await response.json();
					if (errorJson && errorJson.error) {
						errorMsg = errorJson.error;
					}
				} catch {
					console.warn('Failed to parse error response as JSON', response);
				}
				throw new Error(errorMsg);
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
	<Alert>Loading event details...</Alert>
{:else if error}
	<Alert variant="error" title="Error" class="bg-surface">{error}</Alert>
{:else if data}
	<div class="flex flex-col gap-8">
		<section class="rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8">
			<div class="flex flex-col gap-6">
				<SectionHeader eyebrow="Captured event" title={`Event ID: ${data.id}`}>
					{#snippet actions()}
						<Badge variant="neutral" appearance="soft" class="bg-elevated">{data?.method}</Badge>
					{/snippet}
				</SectionHeader>

				<div class="grid gap-4 border-t border-border-muted pt-4 sm:grid-cols-2">
					<div>
						<Eyebrow>Source ID</Eyebrow>
						<p class="mt-2 break-all text-sm leading-6 text-fg">{data.source_id}</p>
					</div>
					<div>
						<Eyebrow>Deduplication hash</Eyebrow>
						<p class="mt-2 break-all text-sm leading-6 text-fg">{data.dedup_hash ?? 'N/A'}</p>
					</div>
				</div>

				<div class="grid gap-4 border-t border-border-muted pt-4 lg:grid-cols-2">
					<section>
						<h4 class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
							Query Parameters
						</h4>
						<KeyValueList
							items={stringArrayRecordToKeyValueItems(data.query_params ?? {})}
							emptyStateText="No values recorded."
						/>
					</section>
					<section>
						<h4 class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">Raw Headers</h4>
						<KeyValueList
							items={stringArrayRecordToKeyValueItems(data.raw_headers ?? {})}
							emptyStateText="No values recorded."
						/>
					</section>
				</div>

				<BodyView body={data.body} contentType={data.body_content_type} />
			</div>
		</section>

		<DeliveryAttemptList {sourceID} {eventID} />
	</div>
{:else}
	<Alert>No details found for this event.</Alert>
{/if}
