<script lang="ts">
	import { page } from '$app/state';
	import EventView from './components/EventView.svelte';

	const sourceID = page.params.sourceID;
	const eventID = page.params.eventID;
</script>

{#if eventID && sourceID}
	<EventView {sourceID} {eventID} />
{:else}
	<p>No source|event ID provided.</p>
{/if}
