<script lang="ts">
	import { resolve } from '$app/paths';
	import Alert from '$lib/components/ui/Alert.svelte';
	import Link from '$lib/components/ui/Link.svelte';
	import { page } from '$app/state';
	import SourceView from './components/SourceView.svelte';

	/** Source id read from the current route. */
	const sourceID = page.params.sourceID;
</script>

<div class="flex flex-col gap-6">
	<nav class="rounded-md border border-border-muted bg-surface px-4 py-3 text-sm shadow-sm">
		<ul class="flex items-center gap-2 text-muted">
			<li><Link href={resolve('/sources')} variant="inline">Sources</Link></li>
			<li>/</li>
			<li class="text-fg">{sourceID || 'Unknown source'}</li>
		</ul>
	</nav>

	<div class="flex flex-col gap-2">
		<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">Source details</p>
		<h1 class="text-3xl font-semibold tracking-tight text-fg">Source ID: {sourceID}</h1>
	</div>

{#if !sourceID}
		<Alert>No source ID provided.</Alert>
{:else}
		<SourceView {sourceID} />
{/if}
</div>
