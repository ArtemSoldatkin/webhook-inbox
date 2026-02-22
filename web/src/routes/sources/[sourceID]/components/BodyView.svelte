<script lang="ts">
	export let body: string;
	export let contentType: string;
</script>

<section>
	<h3>Request body</h3>
	{#if contentType.startsWith('application/json')}
		<pre>{atob(body)}</pre>
	{:else if contentType.startsWith('application/x-www-form-urlencoded')}
		<pre>{atob(body)}</pre>
	{:else if contentType.startsWith('text/plain')}
		<pre>{atob(body)}</pre>
	{:else}
		<p>Unsupported content type: {contentType}</p>
	{/if}
</section>
