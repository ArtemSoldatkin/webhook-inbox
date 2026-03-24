<script lang="ts">
	type Props = {
		/** Group title shown above the rendered entries. */
		title: string;

		/** Key-value pairs where each key maps to multiple values. */
		data: Record<string, string[]>;
	};

	let { title, data }: Props = $props();
</script>

<section>
	<h4>{title}</h4>
	<ul>
		{#each Object.entries(data) as [key, value] (key)}
			<li><strong>{key}:</strong> {value.join(', ')}</li>
		{/each}
	</ul>
</section>
