<script lang="ts">
	type Props = {
		/** Group title shown above the rendered entries. */
		title: string;

		/** Key-value pairs where each key maps to multiple values. */
		data: Record<string, string[]>;
	};

	let { title, data }: Props = $props();
</script>

<section class="rounded-md border border-border-muted bg-surface p-4">
	<h4 class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">{title}</h4>
	{#if Object.keys(data).length > 0}
		<ul class="mt-3 flex flex-col gap-2">
			{#each Object.entries(data) as [key, value] (key)}
				<li
					class="flex flex-col gap-1 rounded-md border border-border-muted bg-elevated px-3 py-2 sm:flex-row sm:items-start sm:justify-between sm:gap-4"
				>
					<span class="text-sm font-medium text-fg">{key}</span>
					<span class="break-all text-sm text-muted">{value.join(', ')}</span>
				</li>
			{/each}
		</ul>
	{:else}
		<p class="mt-2 text-sm text-muted">No values recorded.</p>
	{/if}
</section>
