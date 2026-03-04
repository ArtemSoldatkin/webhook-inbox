<script lang="ts">
	import { type SourceDTO } from '$lib/types';
	export let source: SourceDTO;
</script>

<h2><a href={`/sources/${source.id}`}>{source.id}</a></h2>
<p>{source.description}</p>
<p>{source.ingress_url}</p>
<p>{source.egress_url}</p>
<p>Static headers:</p>
{#each Object.entries(source.static_headers ?? {}) as [key, value]}
	<p>{key}: {value}</p>
{/each}
<p>{source.status}</p>
<p>{source.status_reason}</p>
<p>Created at: {new Date(source.created_at).toLocaleString()}</p>
<p>Updated at: {new Date(source.updated_at).toLocaleString()}</p>
<p>
	Disabled at: {source.disable_at ? new Date(source.disable_at).toLocaleString() : 'N/A'}
</p>
