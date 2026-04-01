<script lang="ts">
	import DOMPurify from 'dompurify';
	import hljs from 'highlight.js/lib/core';
	import 'highlight.js/styles/github.css';
	import formatXML from 'xml-formatter';

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
	<p class="mt-3 text-sm text-warning">{parsed.error}</p>
{/if}
<pre
	class="mt-3 overflow-x-auto rounded-md border border-border-muted bg-elevated p-4 text-sm leading-6 text-fg shadow-sm"><!-- eslint-disable-next-line svelte/no-at-html-tags --><code
		class="hljs xml">{@html parsed.html}</code
	></pre>
