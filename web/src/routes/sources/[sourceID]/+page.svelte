<script lang="ts">
	import { resolve } from '$app/paths';
	import Link from '$lib/components/ui/Link.svelte';
	import { page } from '$app/state';
	import SourceView from './components/SourceView.svelte';

	/** Source id read from the current route. */
	const sourceID = page.params.sourceID;
</script>

<nav>
	<ul>
		<li><Link href={resolve('/sources')} variant="inline">Sources</Link></li>
	</ul>
</nav>
<h2>Source ID: {sourceID}</h2>
{#if !sourceID}
	<p>No source ID provided.</p>
{:else}
	<SourceView {sourceID} />
{/if}
