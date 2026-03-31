<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	/** Type of allowed variants for the alert component. */
	type Variant = 'neutral' | 'info' | 'success' | 'warning' | 'error';

	/** Type of allowed appearances for the alert component. */
	type Appearance = 'soft' | 'outline' | 'solid';

	/** Mapping of appearance styles to appearance. */
	const APPEARANCE_CLASSES: Record<Appearance, Record<Variant, string>> = {
		soft: {
			neutral: 'border-border-muted bg-elevated text-fg',
			success: 'border-success/25 bg-success/10 text-success',
			error: 'border-error/25 bg-error/10 text-error',
			warning: 'border-warning/25 bg-warning/10 text-warning',
			info: 'border-info/25 bg-info/10 text-info'
		},
		outline: {
			neutral: 'border-border bg-transparent text-fg',
			success: 'border-success bg-transparent text-success',
			error: 'border-error bg-transparent text-error',
			warning: 'border-warning bg-transparent text-warning',
			info: 'border-info bg-transparent text-info'
		},
		solid: {
			neutral: 'border-fg bg-fg text-inverted',
			success: 'border-success bg-success text-inverted',
			error: 'border-error bg-error text-inverted',
			warning: 'border-warning bg-warning text-inverted',
			info: 'border-info bg-info text-inverted'
		}
	} as const;

	type Props = Omit<HTMLAttributes<HTMLDivElement>, 'contentEditable' | 'children'> & {
		variant?: Variant;
		appearance?: Appearance;
		title?: string;
		role?: HTMLAttributes<HTMLDivElement>['role'];
		children?: Snippet;
	};

	let {
		variant = 'neutral',
		appearance = 'soft',
		title,
		role = variant === 'error' ? 'alert' : 'status',
		children,
		class: className,
		...rest
	}: Props = $props();
</script>

<div
	class={cx(
		'flex flex-col gap-2 rounded-md border min-h-12 px-4 py-3 text-sm shadow-sm',
		APPEARANCE_CLASSES[appearance][variant],
		className
	)}
	aria-live={variant === 'error' ? 'assertive' : 'polite'}
	{role}
	{...rest}
>
	{#if title}
		<p class="text-sm font-medium">{title}</p>
	{/if}
	{@render children?.()}
</div>
