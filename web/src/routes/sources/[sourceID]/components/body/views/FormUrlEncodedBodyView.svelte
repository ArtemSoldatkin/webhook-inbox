<script lang="ts">
	type Props = {
		/** URL-encoded body content to render. */
		body: string;
	};

	let { body }: Props = $props();

	/** Parsed query parameter pairs from the encoded body. */
	const params = $derived(Array.from(new URLSearchParams(body).entries()));
</script>

{#if params.length > 0}
	<ul class="mt-3 flex flex-col gap-2">
		{#each params as [key, value], index (`${key}:${index}`)}
			<li
				class="flex flex-col gap-1 rounded-md border border-border-muted bg-elevated px-3 py-2 sm:flex-row sm:items-start sm:justify-between sm:gap-4"
			>
				<span class="text-sm font-medium text-fg">{key}</span>
				<span class="break-all text-sm text-muted">{value}</span>
			</li>
		{/each}
	</ul>
{:else}
	<p class="mt-3 text-sm text-muted">No form values provided.</p>
{/if}
