<script lang="ts">
	import DOMPurify from 'dompurify';
	import hljs from 'highlight.js/lib/core';
	import 'highlight.js/styles/github.css';
	import formatXML from 'xml-formatter';

	type Props = {
		body: string;
	};

	let { body }: Props = $props();

	const highlightedXml = $derived(hljs.highlight(formatXML(body), { language: 'xml' }).value);
</script>

<pre><code class="hljs xml">{@html DOMPurify.sanitize(highlightedXml)}</code></pre>
