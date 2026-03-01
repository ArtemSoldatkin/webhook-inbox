<script lang="ts">
	import { type SourceDTO } from '$lib/types';
	export let source: SourceDTO;
</script>

<h2><a href={`/sources/${source.ID}`}>{source.ID}</a></h2>
<p>{source.Description}</p>
<p>{source.IngressUrl}</p>
<p>{source.EgressUrl}</p>
<p>Static headers:</p>
{#each Object.entries(source.StaticHeaders) as [key, value]}
	<p>{key}: {value}</p>
{/each}
<p>{source.Status}</p>
<p>{source.StatusReason}</p>
<p>Created at: {new Date(source.CreatedAt).toLocaleString()}</p>
<p>Updated at: {new Date(source.UpdatedAt).toLocaleString()}</p>
<p>
	Disabled at: {source.DisableAt ? new Date(source.DisableAt).toLocaleString() : 'N/A'}
</p>
