<script lang="ts">
	import hljs from 'highlight.js/lib/core';
	import xml from 'highlight.js/lib/languages/xml';
	import 'highlight.js/styles/github.css';
	import formatXML from 'xml-formatter';

	hljs.registerLanguage('xml', xml);

	export let body: string | undefined;
	export let contentType: string | undefined;

	const parsedBody = body ? atob(body) : '';
	const bodyByteLength = parsedBody.length;
	const highlightedXml = contentType?.startsWith('application/xml')
		? hljs.highlight(formatXML(parsedBody), { language: 'xml' }).value
		: '';

	function downloadBytes() {
		const blob = new Blob([parsedBody], {
			type: contentType || 'application/octet-stream'
		});
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = 'request-body';
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	}
</script>

<section>
	<h3>Request body</h3>
	{#if !body}
		<p>No body content</p>
	{:else if !contentType}
		<p>Content type unknown, cannot display body</p>
	{:else if contentType.startsWith('application/json')}
		<pre>{JSON.stringify(JSON.parse(parsedBody), null, 2)}</pre>
	{:else if contentType.startsWith('application/x-www-form-urlencoded')}
		<pre>{new URLSearchParams(parsedBody).toString()}</pre>
	{:else if contentType.startsWith('application/xml')}
		<pre><code class="hljs xml">{@html highlightedXml}</code></pre>
	{:else}
		<p>Size: {bodyByteLength} bytes</p>
		<button on:click={downloadBytes} disabled={parsedBody === ''}>Download as file</button>
	{/if}
</section>
