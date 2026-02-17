<script lang="ts">
	import { page } from '$app/state';
	import SourceView from './components/SourceView.svelte';

	const { sourceID } = page.params;
</script>

<nav>
	<ul>
		<li><a href="/sources">Sources</a></li>
	</ul>
</nav>
<h2>Source ID: {sourceID}</h2>
{#if !sourceID}
	<p>No source ID provided.</p>
{:else}
	<SourceView {sourceID} />
{/if}
