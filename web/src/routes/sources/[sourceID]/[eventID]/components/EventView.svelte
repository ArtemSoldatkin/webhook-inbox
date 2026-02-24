<script lang="ts">
	import DisplayMapOfStringArrays from '$lib/components/DisplayMapOfStringArrays.svelte';
	import { parseEventDTO } from '$lib/dtoParsers';
	import { type EventDTO } from '$lib/types';
	import BodyView from '../../components/BodyView.svelte';
	import DeliveryAttemptList from './DeliveryAttemptList.svelte';

	export let sourceID: string;
	export let eventID: string;

	let data: EventDTO | null = null;
	let loading = false;
	let error: string | null = null;

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
		<h2>Event ID: {data.ID}</h2>
		<p>Source ID: {data.SourceID}</p>
		<p>Method: {data.Method}</p>
		<DisplayMapOfStringArrays title="Query Parameters" data={data.QueryParams} />
		<DisplayMapOfStringArrays title="Raw Headers" data={data.RawHeaders} />
		<BodyView body={data.Body} contentType={data.BodyContentType} />
	</section>
	<DeliveryAttemptList {sourceID} {eventID} />
{:else}
	<p>No details found for this event.</p>
{/if}
