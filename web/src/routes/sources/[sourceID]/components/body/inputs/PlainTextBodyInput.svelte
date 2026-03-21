<script lang="ts">
	type Props = {
		body: string;
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

	function handleClear() {
		body = '';
		error = null;
	}
</script>

<textarea bind:value={body} rows="10" cols="50" placeholder="Enter body content here..."></textarea>
<button type="button" onclick={handleClear} disabled={!body}>Clear</button>
