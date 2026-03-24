<script lang="ts">
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Props = HTMLButtonAttributes;

	let { type = 'button', children, ...rest }: Props = $props();
</script>

<button {type} {...rest}>
	{@render children?.()}
</button>
