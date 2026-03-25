<script lang="ts">
	type Props = {
		/** An array of objects representing the key-value pairs to be displayed. Each object should have a `label` (the key) and a `value` (the corresponding value).*/
		items: Array<{ label: string; value?: string | number | null }>;

		/** Optional text to display when the items object is empty or undefined. If not provided, no message will be shown in the empty state.*/
		emptyStateText?: string;
	};

	let { items, emptyStateText }: Props = $props();
</script>

{#if items.length > 0}
	<div class="mt-3 flex flex-col gap-2">
		{#each items as { label, value }, index (`${label}-${index}`)}
			<div
				class="flex flex-col gap-1 rounded-md border border-border-muted bg-elevated px-3 py-2 sm:flex-row sm:items-start sm:justify-between sm:gap-4"
			>
				<span class="text-sm font-medium text-fg">{label}</span>
				<span class="break-all text-sm text-muted">{value == null ? 'N/A' : value}</span>
			</div>
		{/each}
	</div>
{:else if emptyStateText}
	<p class="mt-2 text-sm text-muted">{emptyStateText}</p>
{/if}
