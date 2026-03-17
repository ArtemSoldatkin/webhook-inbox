<script lang="ts">
	import type { ContentType } from '$lib/types';

	export let body: string;
	export let contentType: ContentType;

	let error: string | null = null;

	$: {
		switch (contentType) {
			case 'application/json':
				try {
					JSON.parse(body);
					error = null;
				} catch (err: unknown) {
					console.error('Error parsing JSON body:', err);
					error = 'Invalid JSON format';
				}
				break;
			case 'application/x-www-form-urlencoded':
				try {
					new URLSearchParams(body);
					error = null;
				} catch (err: unknown) {
					console.error('Error parsing URL-encoded body:', err);
					error = 'Invalid URL-encoded format';
				}
				break;
			case 'text/plain':
				error = null; // No validation needed for plain text
				break;
			case 'multipart/form-data':
				error = 'Multipart form data is not supported in this editor';
				break;
			case 'application/xml':
				error = null;
				break;
			case 'application/octet-stream':
				error = 'Binary data editing is not supported';
				break;
			default:
				error = null; // No validation for other content types
		}
	}
</script>

<section>
	<select bind:value={contentType}>
		<option value="application/json">JSON</option>
		<option value="application/x-www-form-urlencoded">Form URL Encoded</option>
		<option value="text/plain">Plain Text</option>
		<option value="multipart/form-data">Multipart Form Data</option>
		<option value="application/xml">XML</option>
		<option value="application/octet-stream">Binary Data</option>
	</select>
	<textarea bind:value={body} rows="10" cols="50" placeholder="Enter request body here..."
	></textarea>
	{#if error}
		<div class="error">{error}</div>
	{/if}
</section>
