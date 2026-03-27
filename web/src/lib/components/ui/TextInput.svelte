<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { HTMLInputAttributes } from 'svelte/elements';

	type Props = Omit<HTMLInputAttributes, 'type'>;

	let { value = $bindable(), class: className, ...rest }: Props = $props();
</script>

<input
	type="text"
	bind:value
	class={cx(
		'rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle',
		'focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
		'disabled:pointer-events-none disabled:opacity-50',
		className
	)}
	{...rest}
/>
