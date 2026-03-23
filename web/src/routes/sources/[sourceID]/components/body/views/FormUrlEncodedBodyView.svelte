<script lang="ts">
	type Props = {
		/** URL-encoded body content to render. */
		body: string;
	};

	let { body }: Props = $props();

	/** Parsed query parameter pairs from the encoded body. */
	const params = $derived(Array.from(new URLSearchParams(body).entries()));
</script>

<ul style="list-style-type: none; padding-left: 0">
	{#each params as [key, value], index (`${key}:${index}`)}
		<li>
			<span style="color:#0366d6">{key}</span>: <span style="color:#24292e">{value}</span>
		</li>
	{/each}
</ul>
