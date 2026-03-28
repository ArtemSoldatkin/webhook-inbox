<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { HTMLSelectAttributes } from 'svelte/elements';

	/** A single select option. */
	type SelectOption = {
		/** Submitted value. */
		value: string | number;

		/** Visible label. */
		label: string;

		/** Whether the option is disabled. */
		disabled?: boolean;
	};

	type Props = Omit<HTMLSelectAttributes, 'value' | 'multiple'> & {
		/** Selected value. */
		value?: string | number;

		/** Available options. */
		options: SelectOption[];
	};

	let { value = $bindable(), options, class: className, ...rest }: Props = $props();
</script>

<select
	bind:value
	class={cx(
		'rounded-md border border-border bg-elevated min-h-12 px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle',
		'focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
		'disabled:pointer-events-none disabled:opacity-50',
		className
	)}
	{...rest}
>
	{#each options as option (option.value)}
		<option value={option.value} disabled={option.disabled}>
			{option.label}
		</option>
	{/each}
</select>
