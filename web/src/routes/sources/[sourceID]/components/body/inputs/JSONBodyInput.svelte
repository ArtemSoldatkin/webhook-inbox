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

<textarea bind:value={body} rows="10" cols="50" placeholder="Enter JSON body here..."></textarea>
<Button type="button" onclick={formatJSON} disabled={!!error}>Format</Button>
<Button type="button" onclick={handleClear} disabled={!body}>Clear</Button>
