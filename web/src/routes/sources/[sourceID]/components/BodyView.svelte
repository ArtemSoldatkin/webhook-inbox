<script lang="ts">
	export let body: string | undefined;
	export let contentType: string | undefined;
</script>

<section>
	<h3>Request body</h3>
	{#if !body}
		<p>No body content</p>
	{:else if !contentType}
		<p>Content type unknown, cannot display body</p>
	{:else if contentType.startsWith('application/json')}
		<pre>{JSON.stringify(JSON.parse(atob(body)), null, 2)}</pre>
	{:else if contentType.startsWith('application/x-www-form-urlencoded')}
		<pre>{new URLSearchParams(atob(body)).toString()}</pre>
	{:else if contentType.startsWith('text/plain')}
		<pre>{atob(body)}</pre>
	{:else}
		<p>Unsupported content type: {contentType}</p>
	{/if}
</section>
