<script lang="ts">
	import DisplayMapOfStringArrays from '$lib/components/DisplayMapOfStringArrays.svelte';
	import { parseEventDTO } from '$lib/dto-parsers';
	import { type EventDTO } from '$lib/types';
	import BodyView from '../../components/BodyView.svelte';
	import DeliveryAttemptList from './DeliveryAttemptList.svelte';

	/** Source id owning the current event. */
	export let sourceID: string;

	/** Event id shown on the details page. */
	export let eventID: string;

	/** Loaded event details. */
	let data: EventDTO | null = null;

	/** Tracks whether the event request is in flight. */
	let loading = false;

	/** Holds the latest event loading error. */
	let error: string | null = null;

	/** Fetches the current event and normalizes its payload. */
	async function fetchEventDetails() {
		loading = true;
		error = null;
		try {
			const response = await fetch(`/api/sources/${sourceID}/events/${eventID}`);
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

	$: if (eventID) {
		fetchEventDetails();
	}
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
