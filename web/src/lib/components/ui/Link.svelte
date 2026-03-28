	<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { Snippet } from 'svelte';
	import type { HTMLAnchorAttributes } from 'svelte/elements';

	/**
	 * Variant of the link, determining its styling.
	 * - `primary`: The default link style, typically used for main actions.
	 * - `secondary`: A less prominent style, often used for secondary actions.
	 * - `inverse`: A style suitable for dark backgrounds, with light text and borders.
	 * - `nav`: A style designed for navigation links, often with specific hover and active states.
	 * - `brand`: A style that incorporates brand colors and design elements.
	 * - `inline`: A style that makes the link appear inline with surrounding text, often without additional padding or background.
	 */
	type Variant = 'primary' | 'secondary' | 'inverse' | 'nav' | 'brand' | 'inline';

	/**
	 * Mapping of link variants to their corresponding Tailwind CSS classes.
	 * Each variant has a specific set of classes that define its appearance and behavior,
	 * including background color, text color, shadow, hover effects, active state, focus state, and disabled state.
	 * This mapping allows for easy application of styles based on the selected variant.
	 */
	const VARIANT_CLASSES: Record<Variant, string> = {
		primary:
			'inline-flex items-center justify-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium whitespace-nowrap text-inverted shadow-sm hover:bg-primary-hover hover:shadow-md active:bg-primary-active focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
		secondary:
			'inline-flex items-center justify-center gap-2 rounded-md border border-border bg-surface px-4 py-2 text-sm font-medium whitespace-nowrap text-fg shadow-sm hover:bg-elevated hover:shadow-md active:bg-bg focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
		inverse:
			'inline-flex items-center justify-center gap-2 rounded-md border border-inverted/20 bg-transparent px-4 py-2 text-sm font-medium whitespace-nowrap text-inverted hover:bg-primary-hover active:bg-primary-active focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
		nav: 'inline-flex items-center rounded-md px-3 py-2 text-sm font-medium text-muted hover:bg-elevated hover:text-fg focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
		brand:
			'inline-flex items-center rounded-md text-base font-semibold tracking-tight text-fg hover:text-primary focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
		inline:
			'text-primary underline-offset-4 hover:underline focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2'
	} as const;

	type Props = Omit<HTMLAnchorAttributes, 'children'> & {
		variant?: Variant;
		children?: Snippet;
	};

	let { class: className, variant = 'primary', children, ...rest }: Props = $props();
</script>

<a class={cx('transition-colors outline-none disabled:pointer-events-none disabled:opacity-50', VARIANT_CLASSES[variant], className)} {...rest}>
	{@render children?.()}
</a>
