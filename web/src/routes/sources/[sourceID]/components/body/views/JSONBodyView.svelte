<script lang="ts">
	import DOMPurify from 'dompurify';
	import hljs from 'highlight.js/lib/core';
	import 'highlight.js/styles/github.css';

	type Props = {
		/** JSON body content to render. */
		body: string;
	};

	let { body }: Props = $props();

	/** Highlighted JSON output and any parse error. */
	const parsed = $derived.by(() => {
		try {
			const formatted = JSON.stringify(JSON.parse(body), null, 2);
			const html = hljs.highlight(formatted, { language: 'json' }).value;

			return {
				html: DOMPurify.sanitize(html),
				error: null
			};
		} catch {
			return {
				html: DOMPurify.sanitize(hljs.highlight(body, { language: 'json' }).value),
				error: 'Invalid JSON'
			};
		}
	});
</script>

{#if parsed.error}
	<p class="mt-3 text-sm text-warning">{parsed.error}</p>
{/if}
<!-- eslint-disable-next-line svelte/no-at-html-tags -->
<pre class="mt-3 overflow-x-auto rounded-md border border-border-muted bg-elevated p-4 text-sm leading-6 text-fg shadow-sm"><code class="hljs json">{@html parsed.html}</code></pre>
