<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';

	type Props = {
		/** Bound JSON body value. */
		body: string;

		/** Validation error shown by the input. */
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

	/** Pretty-prints the current JSON body when it is valid. */
	function formatJSON(): void {
		try {
			const parsed = JSON.parse(body);
			body = JSON.stringify(parsed, null, 2);
		} catch (err) {
			console.error('Error formatting JSON body:', err);
			error = 'Cannot format invalid JSON';
		}
	}

	/** Clears the current JSON body. */
	function handleClear(): void {
		body = '';
		error = null;
	}

	$effect(() => {
		try {
			JSON.parse(body);
			error = null;
		} catch {
			error = 'Invalid JSON format';
		}
	});
</script>

<div class="flex flex-col gap-4">
	<textarea
		bind:value={body}
		rows="10"
		placeholder="Enter JSON body here..."
		class="min-h-56 w-full rounded-md border border-border bg-surface px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
	></textarea>
	<div class="flex flex-wrap justify-end gap-3">
		<Button type="button" onclick={formatJSON} disabled={!!error} variant="secondary">Format</Button>
		<Button type="button" onclick={handleClear} disabled={!body} variant="secondary">Clear</Button>
	</div>
</div>
