<script lang="ts">
	import DOMPurify from 'dompurify';
	import hljs from 'highlight.js/lib/core';
	import 'highlight.js/styles/github.css';

	type Props = {
		body: string;
	};

	let { body }: Props = $props();

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
	<p>{parsed.error}</p>
{/if}
<pre><code class="hljs xml">{@html parsed.html}</code></pre>
