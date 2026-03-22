<script lang="ts">
	import DOMPurify from 'dompurify';
	import hljs from 'highlight.js/lib/core';
	import 'highlight.js/styles/github.css';
	import formatXML from 'xml-formatter';

	/** Props for rendering XML request bodies. */
	type Props = {
		/** XML body content to render. */
		body: string;
	};

	let { body }: Props = $props();

	/** Highlighted XML output and any parse error. */
	const parsed = $derived.by(() => {
		try {
			const formatted = formatXML(body);
			const html = hljs.highlight(formatted, { language: 'xml' }).value;
			return {
				html: DOMPurify.sanitize(html),
				error: null
			};
		} catch {
			return {
				html: DOMPurify.sanitize(hljs.highlight(body, { language: 'xml' }).value),
				error: 'Invalid XML'
			};
		}
	});
</script>

{#if parsed.error}
	<p>{parsed.error}</p>
{/if}
<pre><code class="hljs xml">{@html parsed.html}</code></pre>
