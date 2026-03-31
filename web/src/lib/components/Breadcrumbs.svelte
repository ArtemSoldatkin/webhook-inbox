<script lang="ts">
	import Link from './ui/Link.svelte';

	/** Type definition for a single breadcrumb item. */
	type BreadcrumbItem = {
		/** The display label for the breadcrumb item. */
		label: string;
		/** Optional URL to link to when the breadcrumb item is clicked.
		 * If not provided, the item will be rendered as plain text. */
		href?: string;
		/** Optional flag to indicate if this breadcrumb item represents the current active page.
		 * If true, the item will be rendered as plain text regardless of the presence of href. */
		active?: boolean;
	};

	type Props = {
		/** List of breadcrumb items to display. */
		items: BreadcrumbItem[];
	};

	let { items }: Props = $props();
</script>

<nav class="rounded-md border border-border-muted bg-surface px-4 py-3 text-sm shadow-sm">
	<ul class="flex flex-wrap items-center gap-2 text-muted">
		{#each items as item, index (index)}
			<li>
				{#if item.href && !item.active}
					<Link href={item.href} variant="inline">{item.label}</Link>
				{:else}
					<span class="text-fg">{item.label}</span>
				{/if}
			</li>
			{#if index < items.length - 1}
				<li>/</li>
			{/if}
		{/each}
	</ul>
</nav>
