<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		/** The button type attribute, which can be 'button', 'submit', or 'reset'. Defaults to 'button'. */
		type?: 'button' | 'submit' | 'reset';

		/** A boolean indicating whether the button is disabled. Defaults to false. */
		disabled?: boolean;

		/** The content to be rendered inside the button, which can be any valid Svelte snippet. */
		children: Snippet<[]>;

		/** An optional click event handler function that will be called when the button is clicked. */
		onClick?: () => void;
	};

	let { type = 'button', disabled = false, children, onClick }: Props = $props();
</script>

<button {type} {disabled} onclick={onClick}>
	{@render children()}
</button>
