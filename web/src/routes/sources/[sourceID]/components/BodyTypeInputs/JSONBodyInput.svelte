<script lang="ts">
	type Props = {
		body: string;
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

	function formatJSON() {
		try {
			const parsed = JSON.parse(body);
			body = JSON.stringify(parsed, null, 2);
		} catch (err) {
			console.error('Error formatting JSON body:', err);
			error = 'Cannot format invalid JSON';
		}
	}

	function handleClear() {
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
<button type="button" onclick={formatJSON} disabled={!!error}>Format</button>
<button type="button" onclick={handleClear}>Clear</button>
