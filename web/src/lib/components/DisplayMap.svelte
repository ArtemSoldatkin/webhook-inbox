<script lang="ts">
	export let data: Record<string, string>;
</script>

<div>
	{#each Object.entries(data) as [key, value]}
		<p><strong>{key}:</strong> {value}</p>
	{/each}
</div>
