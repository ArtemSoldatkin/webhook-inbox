	<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	/** Type of allowed variants for the Eyebrow component. */
	type Variant = 'strong' | 'subtle';

	/** Mapping of variant styles to variant. */
	const VARIANT_CLASSES: Record<Variant, string> = {
		strong: 'text-sm tracking-[0.18em] text-primary',
		subtle: 'text-xs tracking-[0.12em] text-subtle'
	} as const;

	type Props =
		| (Omit<HTMLAttributes<HTMLParagraphElement>, 'contentEditable' | 'children'> & {
				as?: 'p';
				variant?: Variant;
				children?: Snippet;
		  })
		| (Omit<HTMLAttributes<HTMLLabelElement>, 'children'> & {
				as: 'label';
				variant?: Variant;
				children?: Snippet;
		  });

	let { as = 'p', variant = 'subtle', class: className, children, ...rest }: Props = $props();
</script>

<svelte:element
	this={as}
	class={cx('font-medium uppercase', VARIANT_CLASSES[variant], className)}
	{...rest}
>
	{@render children?.()}
</svelte:element>
