<script lang="ts">
	import { cx } from '$lib/utils/cx';
	import type { HTMLInputAttributes } from 'svelte/elements';

	type Props = Omit<HTMLInputAttributes, 'type' | 'value' | 'on:change' | 'onchange'> & {
		/** Bound numeric value of the input. */
		value?: number | null;

		/** Optional change handler that receives the input change event. */
		onchange?: (event: Event) => void;
	};

	let { value = $bindable(), class: className, onchange, ...rest }: Props = $props();

	/**
	 * Handles changes to the input, parsing the value as a float
	 * and updating the bound value if valid. Resets to null if the input is cleared.
	 *
	 * @param event - The input change event containing the new value.
	 */
	function handleChange(event: Event): void {
		const targetValue = (event.target as HTMLInputElement).value;
		const parsedValue = parseFloat(targetValue);

		if (!isNaN(parsedValue)) {
			value = parsedValue;
		} else if (targetValue === '') {
			value = null;
		}

		onchange?.(event); // Call any additional onchange handler passed via props
	}
</script>

<input
	type="number"
	value={value == null ? '' : value.toString()}
	class={cx(
		'rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle',
		'focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
		'disabled:pointer-events-none disabled:opacity-50',
		className
	)}
	onchange={handleChange}
	{...rest}
/>
