<script lang="ts">
	import { type SourceDTO } from '$lib/types';

	import DisplayPaginatedList from '$lib/components/DisplayPaginatedList.svelte';
	import CreateSource from './components/CreateSource.svelte';
	import SourceView from './components/SourceView.svelte';
</script>

<CreateSource />
<!-- <ListSources /> -->
<DisplayPaginatedList url="/api/sources">
	<section let:item>
		<SourceView source={item as SourceDTO} />
	</section>
</DisplayPaginatedList>
