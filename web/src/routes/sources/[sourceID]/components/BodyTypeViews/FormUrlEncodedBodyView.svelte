<script lang="ts">
	type Props = {
		body: string;
	};

	let { body }: Props = $props();
</script>

<ul style="list-style-type: none; padding-left: 0">
	{#each Array.from(new URLSearchParams(body)) as [key, value] (key)}
		<li>
			<span style="color:#0366d6">{key}</span>: <span style="color:#24292e">{value}</span>
		</li>
	{/each}
</ul>
