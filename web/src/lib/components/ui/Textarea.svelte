<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { HTMLTextareaAttributes } from 'svelte/elements';

	type Props = HTMLTextareaAttributes & {
		/** Bound value of the textarea. */
		value?: string;
	};

	let { value = $bindable(), class: className, rows = 10, ...rest }: Props = $props();
</script>

<textarea
	{rows}
	bind:value
	class={cx(
		'min-h-56 w-full rounded-md border border-border bg-surface px-4 py-3',
		'text-sm text-fg shadow-sm outline-none',
		'placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
		'disabled:pointer-events-none disabled:opacity-50',
		className
	)}
	{...rest}
></textarea>
