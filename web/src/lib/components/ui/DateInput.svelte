<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { HTMLInputAttributes } from 'svelte/elements';

	type Props = Omit<HTMLInputAttributes, 'type' | 'value'> & {
		/** Bound string value of the date input, representing the selected date in YYYY-MM-DD format. */
		value?: string | null;
	};

	let { value = $bindable(), class: className, ...rest }: Props = $props();
</script>

<input
	type="date"
	bind:value
	class={cx(
		'rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
		className
	)}
	{...rest}
/>
