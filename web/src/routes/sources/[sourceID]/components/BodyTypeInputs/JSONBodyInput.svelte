<script lang="ts">
	type Props = {
		body: string;
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

	$effect(() => {
		try {
			JSON.parse(body);
			error = null;
		} catch (err: unknown) {
			console.error('Error parsing JSON body:', err);
			error = 'Invalid JSON format';
		}
	});
</script>

<textarea bind:value={body} rows="10" cols="50" placeholder="Enter request body here..."></textarea>
