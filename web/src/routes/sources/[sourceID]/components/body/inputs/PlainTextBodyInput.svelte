<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';

	type Props = {
		/** Bound plain text body value. */
		body: string;

		/** Validation error shown by the input. */
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

	/** Clears the current plain text body. */
	function handleClear(): void {
		body = '';
		error = null;
	}
</script>

<textarea bind:value={body} rows="10" cols="50" placeholder="Enter body content here..."></textarea>
<Button type="button" onclick={handleClear} disabled={!body}>Clear</Button>
