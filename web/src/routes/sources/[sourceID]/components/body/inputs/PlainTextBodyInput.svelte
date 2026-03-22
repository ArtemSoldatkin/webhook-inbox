<script lang="ts">
	/** Props shared by plain text body input controls. */
	type Props = {
		/** Bound plain text body value. */
		body: string;

		/** Validation error shown by the input. */
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

	/** Clears the current plain text body. */
	function handleClear() {
		body = '';
		error = null;
	}
</script>

<textarea bind:value={body} rows="10" cols="50" placeholder="Enter body content here..."></textarea>
<button type="button" onclick={handleClear} disabled={!body}>Clear</button>
