<script lang="ts">
	import DOMPurify from 'dompurify';
	import hljs from 'highlight.js/lib/core';
	import 'highlight.js/styles/github.css';

	type Props = {
		body: string;
	};

	let { body }: Props = $props();

	const highlightedXml = $derived(
		hljs.highlight(JSON.stringify(JSON.parse(body), null, 2), {
			language: 'json'
		}).value
	);
</script>

<pre><code class="hljs xml">{@html DOMPurify.sanitize(highlightedXml)}</code></pre>
