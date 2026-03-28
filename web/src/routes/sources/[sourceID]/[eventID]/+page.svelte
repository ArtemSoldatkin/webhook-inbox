<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Alert from '$lib/components/ui/Alert.svelte';
	import Link from '$lib/components/ui/Link.svelte';
	import SectionHeader from '$lib/components/ui/SectionHeader.svelte';
	import EventView from './components/EventView.svelte';

	/** Source id read from the current route. */
	const sourceID = page.params.sourceID;

	/** Event id read from the current route. */
	const eventID = page.params.eventID;
</script>

<div class="flex flex-col gap-6">
	<nav class="rounded-md border border-border-muted bg-surface px-4 py-3 text-sm shadow-sm">
		<ul class="flex flex-wrap items-center gap-2 text-muted">
			<li><Link href={resolve('/sources')} variant="inline">Sources</Link></li>
			<li>/</li>
			<li>
				<Link href={resolve(`/sources/${sourceID}`)} variant="inline">{sourceID}</Link>
			</li>
			<li>/</li>
			<li class="text-fg">{eventID || 'Unknown event'}</li>
		</ul>
	</nav>

	<SectionHeader eyebrow="Event details" title={`Event ID: ${eventID}`} titleAs="h1" />

	{#if eventID && sourceID}
		<EventView {sourceID} {eventID} />
	{:else}
		<Alert>No source/event ID provided.</Alert>
	{/if}
</div>
