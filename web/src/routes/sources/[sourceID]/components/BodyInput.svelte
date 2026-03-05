<script lang="ts">
	import type { ContentType } from '$lib/types';

	export let body: string;
	export let contentType: ContentType;

	let error: string | null = null;

	$: {
		if (contentType === 'application/json') {
			try {
				JSON.parse(body);
				error = null;
			} catch (err: unknown) {
				console.error('Error parsing JSON body:', err);
				error = 'Invalid JSON format';
			}
		} else if (contentType === 'application/x-www-form-urlencoded') {
			try {
				new URLSearchParams(body);
				error = null;
			} catch (err: unknown) {
				console.error('Error parsing URL-encoded body:', err);
				error = 'Invalid URL-encoded format';
			}
		} else {
			error = null; // No validation for plain text
		}
	}
</script>

<section>
	<select bind:value={contentType}>
		<option value="application/json">JSON</option>
		<option value="application/x-www-form-urlencoded">Form URL Encoded</option>
		<option value="text/plain">Plain Text</option>
	</select>
	<textarea bind:value={body} rows="10" cols="50" placeholder="Enter request body here..."
	></textarea>
	{#if error}
		<div class="error">{error}</div>
	{/if}
</section>
