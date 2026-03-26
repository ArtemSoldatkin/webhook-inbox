<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { HTMLAttributes } from 'svelte/elements';

	/** Type of allowed variants for the SectionEyebrow component. */
	type Variant = 'strong' | 'subtle';

	/** Mapping of variant styles to variant. */
	const VARIANT_CLASSES: Record<Variant, string> = {
		strong: 'text-sm tracking-[0.18em] text-primary',
		subtle: 'text-xs tracking-[0.12em] text-subtle'
	} as const;

	type Props = Omit<HTMLAttributes<HTMLParagraphElement>, 'contentEditable'> & {
		variant?: Variant;
	};

	let { variant = 'subtle', class: className, children, ...rest }: Props = $props();
</script>

<p class={cx('font-medium uppercase', VARIANT_CLASSES[variant], className)} {...rest}>
	{@render children?.()}
</p>
