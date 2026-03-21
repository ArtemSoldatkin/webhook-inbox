<script lang="ts">
	import InputMap from '$lib/components/InputMap.svelte';

	type Props = {
		body: string;
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();
	let jsonBody = $state<Record<string, string>>({});

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
