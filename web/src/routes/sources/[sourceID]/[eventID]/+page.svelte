<script lang="ts">
	import { page } from '$app/state';
	import EventView from './components/EventView.svelte';

	const { sourceID, eventID } = page.params;
</script>

{#if eventID && sourceID}
	<EventView {sourceID} {eventID} />
{:else}
	<p>No source|event ID provided.</p>
{/if}
