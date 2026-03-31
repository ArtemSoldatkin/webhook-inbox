<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Textarea from '$lib/components/ui/Textarea.svelte';

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
	<Textarea
		aria-label="Plain text request body"
		bind:value={body}
		rows={10}
		placeholder="Enter body content here..."
	/>
	<div class="flex justify-end">
		<Button type="button" onclick={handleClear} disabled={!body} variant="secondary">Clear</Button>
	</div>
</div>
