<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Alert from '$lib/components/ui/Alert.svelte';
	import Link from '$lib/components/ui/Link.svelte';
	import SectionHeader from '$lib/components/ui/SectionHeader.svelte';
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

	<SectionHeader eyebrow="Source details" title={`Source ID: ${sourceID}`} titleAs="h1" />

	{#if !sourceID}
		<Alert>No source ID provided.</Alert>
	{:else}
		<SourceView {sourceID} />
	{/if}
</div>
