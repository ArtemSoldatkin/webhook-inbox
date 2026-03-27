	<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	/** Type of allowed variants for the badge component. */
	type Variant = 'neutral' | 'success' | 'error' | 'warning' | 'info';

	/** Type of allowed appearances for the badge component. */
	type Appearance = 'solid' | 'soft' | 'outline';

	/** Mapping of appearance styles to appearance. */
	const APPEARANCE_CLASSES: Record<Appearance, Record<Variant, string>> = {
		soft: {
			neutral: 'border-border bg-surface text-muted',
			success: 'border-success/25 bg-success/10 text-success',
			error: 'border-error/25 bg-error/10 text-error',
			warning: 'border-warning/25 bg-warning/10 text-warning',
			info: 'border-info/25 bg-info/10 text-info'
		},
		outline: {
			neutral: 'border-border bg-transparent text-muted',
			success: 'border-success bg-transparent text-success',
			error: 'border-error bg-transparent text-error',
			warning: 'border-warning bg-transparent text-warning',
			info: 'border-info bg-transparent text-info'
		},
		solid: {
			neutral: 'border-border bg-fg text-inverted',
			success: 'border-success bg-success text-inverted',
			error: 'border-error bg-error text-inverted',
			warning: 'border-warning bg-warning text-inverted',
			info: 'border-info bg-info text-inverted'
		}
	} as const;

	type Props =
		| (Omit<HTMLAttributes<HTMLSpanElement>, 'contentEditable' | 'children'> & {
				variant?: Variant;
				appearance?: Appearance;
				as?: 'span';
				children?: Snippet;
		  })
		| (Omit<HTMLAttributes<HTMLDivElement>, 'contentEditable' | 'children'> & {
				variant?: Variant;
				appearance?: Appearance;
				as: 'div';
				children?: Snippet;
		  });

	let {
		variant = 'neutral',
		appearance = 'solid',
		as = 'span',
		class: className,
		children,
		...rest
	}: Props = $props();
</script>

<svelte:element
	this={as}
	class={cx(
		'inline-flex w-fit items-center rounded-full border px-3 py-1 text-xs font-medium uppercase tracking-[0.12em]',
		APPEARANCE_CLASSES[appearance][variant],
		className
	)}
	{...rest}
>
	{@render children?.()}
</svelte:element>
