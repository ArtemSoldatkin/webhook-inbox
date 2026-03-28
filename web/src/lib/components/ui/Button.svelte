<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	/**
	 * Variant of the button, determining its styling.
	 * - `primary`: The default button style, used for main actions.
	 * - `secondary`: A less prominent style, used for secondary actions.
	 * - `ghost`: A minimal style with no background, used for subtle actions.
	 * - `destructive`: A style indicating a destructive action, such as delete.
	 * - `link`: A style that looks like a hyperlink, used for navigation or less important actions.
	 */
	type Variant = 'primary' | 'secondary' | 'ghost' | 'destructive' | 'link';

	/**
	 * Mapping of button variants to their corresponding Tailwind CSS classes.
	 * Each variant has a specific set of classes that define its appearance and behavior,
	 * including background color, text color, shadow, hover effects, active state, focus state, and disabled state.
	 * This mapping allows for easy application of styles based on the selected variant.
	 */
	const VARIANT_CLASSES: Record<Variant, string> = {
		primary:
			'bg-primary text-inverted shadow-sm hover:bg-primary-hover hover:shadow-md active:bg-primary-active focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
		secondary:
			'bg-surface text-fg shadow-sm hover:bg-elevated hover:shadow-md active:bg-bg focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
		ghost:
			'bg-transparent text-primary hover:bg-gray-100 active:bg-gray-200 focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
		destructive:
			'bg-error text-inverted shadow-sm hover:bg-error hover:shadow-md active:bg-error focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
		link: 'bg-transparent text-primary underline-offset-4 hover:underline focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50'
	} as const;

	type Props = Omit<HTMLButtonAttributes, 'children' | 'type'> & {
		variant?: Variant;
		type?: HTMLButtonAttributes['type'];
		children?: Snippet;
	};

	let {
		class: className,
		variant = 'primary',
		type = 'button',
		children,
		...rest
	}: Props = $props();
</script>

<button
	class={cx(
		'inline-flex items-center justify-center gap-2 rounded-md min-h-12 px-4 py-2 text-sm font-medium whitespace-nowrap',
		'transition-colors outline-none',
		VARIANT_CLASSES[variant],
		className
	)}
	{type}
	{...rest}
>
	{@render children?.()}
</button>
