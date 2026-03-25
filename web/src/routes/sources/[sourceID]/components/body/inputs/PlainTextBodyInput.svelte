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

<div class="flex flex-col gap-4">
	<textarea
		bind:value={body}
		rows="10"
		placeholder="Enter body content here..."
		class="min-h-56 w-full rounded-md border border-border bg-surface px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
	></textarea>
	<div class="flex justify-end">
		<Button type="button" onclick={handleClear} disabled={!body} variant="secondary">Clear</Button>
	</div>
</div>
