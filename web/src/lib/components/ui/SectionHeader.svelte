<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
	import Eyebrow from './Eyebrow.svelte';

	type TitleAs = 'h1' | 'h2' | 'h3';

	const TITLE_CLASSES: Record<TitleAs, string> = {
		h1: 'mt-4 text-3xl font-semibold tracking-tight text-fg',
		h2: 'mt-4 text-3xl font-semibold tracking-tight text-fg',
		h3: 'mt-4 text-2xl font-semibold tracking-tight text-fg'
	};

	type Props = Omit<HTMLAttributes<HTMLDivElement>, 'children'> & {
		/** Optional eyebrow shown above the title. */
		eyebrow?: string;

		/** Main section title. */
		title: string;

		/** Optional helper text shown under the title. */
		description?: string;

		/** Heading element used for the title. */
		titleAs?: TitleAs;

		/** Optional actions shown beside the title block. */
		actions?: Snippet;
	};

	let {
		eyebrow,
		title,
		description,
		titleAs = 'h2',
		actions,
		class: className,
		...rest
	}: Props = $props();
</script>

<div
	class={cx('flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between', className)}
	{...rest}
>
	<div class="min-w-0 flex-1">
		{#if eyebrow}
			<Eyebrow variant="strong">{eyebrow}</Eyebrow>
		{/if}

		<svelte:element this={titleAs} class={TITLE_CLASSES[titleAs]}>
			{title}
		</svelte:element>

		{#if description}
			<p class="mt-3 max-w-2xl text-sm leading-6 text-muted sm:text-base">
				{description}
			</p>
		{/if}
	</div>

	{#if actions}
		<div class="shrink-0">
			{@render actions()}
		</div>
	{/if}
</div>
