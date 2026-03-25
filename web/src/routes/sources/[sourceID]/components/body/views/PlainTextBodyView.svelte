<script lang="ts">
	type Props = {
		/** Plain text body content to render. */
		body: string;
	};

	/** Plain text body content to display. */
	let { body }: Props = $props();
</script>

<pre class="mt-3 overflow-x-auto rounded-md border border-border-muted bg-elevated p-4 text-sm leading-6 text-fg shadow-sm">{body}</pre>
