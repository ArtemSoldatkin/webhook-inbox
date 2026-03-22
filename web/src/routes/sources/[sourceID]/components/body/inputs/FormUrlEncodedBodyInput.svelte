<script lang="ts">
	import InputMap from '$lib/components/InputMap.svelte';

	/** Props shared by URL-encoded body input controls. */
	type Props = {
		/** Bound encoded body value. */
		body: string;

		/** Validation error shown by the input. */
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

	/** Editable key-value entries before serialization. */
	let jsonBody = $state<Record<string, string>>({});

	/** Clears the current URL-encoded body. */
	function handleClear() {
		body = '';
		error = null;
	}

	$effect(() => {
		const params = Object.entries(jsonBody);
		if (params.length > 0) {
			body = params
				.map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
				.join('&');
		} else {
			body = '';
		}
	});
</script>

<InputMap bind:map={jsonBody} disabled={!!error} />
<button type="button" onclick={handleClear} disabled={!body}>Clear</button>
