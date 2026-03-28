<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { HTMLInputAttributes } from 'svelte/elements';

	type Props = Omit<HTMLInputAttributes, 'type' | 'checked' | 'value'> & {
		/** Bound boolean value of the checkbox. */
		value?: boolean;
	};

	let { value = $bindable(), class: className, ...rest }: Props = $props();
</script>

<input
	type="checkbox"
	class={cx(
		'rounded border border-border bg-elevated text-fg shadow-sm outline-none',
		'focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
		'disabled:pointer-events-none disabled:opacity-50',
		className
	)}
	bind:checked={value}
	{...rest}
/>
