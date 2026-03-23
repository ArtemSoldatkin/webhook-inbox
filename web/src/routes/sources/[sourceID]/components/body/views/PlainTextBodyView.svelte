<script lang="ts">
	type Props = {
		/** Plain text body content to render. */
		body: string;
	};

	/** Plain text body content to display. */
	let { body }: Props = $props();
</script>

<pre>{body}</pre>
