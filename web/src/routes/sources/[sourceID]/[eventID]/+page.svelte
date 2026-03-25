<script lang="ts">
	import { resolve } from '$app/paths';
	import Link from '$lib/components/ui/Link.svelte';
	import { page } from '$app/state';
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

	<div class="flex flex-col gap-2">
		<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">Event details</p>
		<h1 class="text-3xl font-semibold tracking-tight text-fg">Event ID: {eventID}</h1>
	</div>

	{#if eventID && sourceID}
		<EventView {sourceID} {eventID} />
	{:else}
		<div class="rounded-md border border-border-muted bg-elevated px-4 py-6 text-sm text-muted">
			No source/event ID provided.
		</div>
	{/if}
</div>
