<script lang="ts">
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
	<p>Loading event details...</p>
{:else if error}
	<p class="error">Error: {error}</p>
{:else if data}
	<section>
		<h2>Event ID: {data.id}</h2>
		<p>Source ID: {data.source_id}</p>
		<p>Deduplication Hash: {data.dedup_hash ?? 'N/A'}</p>
		<p>Method: {data.method}</p>
		<DisplayMapOfStringArrays title="Query Parameters" data={data.query_params ?? {}} />
		<DisplayMapOfStringArrays title="Raw Headers" data={data.raw_headers ?? {}} />
		<BodyView body={data.body} contentType={data.body_content_type} />
	</section>
	<DeliveryAttemptList {sourceID} {eventID} />
{:else}
	<p>No details found for this event.</p>
{/if}
