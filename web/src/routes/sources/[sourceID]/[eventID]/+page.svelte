<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Breadcrumbs from '$lib/components/Breadcrumbs.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import SectionHeader from '$lib/components/ui/SectionHeader.svelte';
	import EventView from './components/EventView.svelte';

	/** Source id read from the current route. */
	const sourceID = page.params.sourceID;

	/** Event id read from the current route. */
	const eventID = page.params.eventID;

	/** Breadcrumb items for navigation, dynamically generated based on the current source and event IDs. */
	const breadcrumbItems = [
		{ label: 'Sources', href: resolve('/sources') },
		{ label: sourceID ?? 'Unknown source', href: resolve(`/sources/${sourceID}`) },
		{ label: eventID ?? 'Unknown event', active: true }
	];
</script>

<div class="flex flex-col gap-6">
	<Breadcrumbs items={breadcrumbItems} />

	<SectionHeader eyebrow="Event details" title={`Event ID: ${eventID}`} titleAs="h1" />

	{#if eventID && sourceID}
		<EventView {sourceID} {eventID} />
	{:else}
		<Alert>No source/event ID provided.</Alert>
	{/if}
</div>
