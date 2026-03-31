<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Breadcrumbs from '$lib/components/Breadcrumbs.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import SectionHeader from '$lib/components/ui/SectionHeader.svelte';
	import SourceView from './components/SourceView.svelte';

	/** Source id read from the current route. */
	const sourceID = page.params.sourceID;

	/** Breadcrumb items for navigation, dynamically generated based on the current source ID. */
	const breadcrumbItems = [
		{ label: 'Sources', href: resolve('/sources') },
		{ label: sourceID ?? 'Unknown source', active: true }
	];
</script>

<div class="flex flex-col gap-6">
	<Breadcrumbs items={breadcrumbItems} />

	<SectionHeader eyebrow="Source details" title={`Source ID: ${sourceID}`} titleAs="h1" />

	{#if !sourceID}
		<Alert>No source ID provided.</Alert>
	{:else}
		<SourceView {sourceID} />
	{/if}
</div>
