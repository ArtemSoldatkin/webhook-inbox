<script lang="ts">
	import type { ContentType } from '$lib/types';

	type Props = {
		body: string;
		contentType?: ContentType;
	};

	let { body, contentType }: Props = $props();

	const bodyByteLength = $derived(body.length);
	const hexPreview = $derived(
		Array.from(body.slice(0, 16))
			.map((c) => c.charCodeAt(0).toString(16).padStart(2, '0'))
			.join(' ')
	);

	function downloadBytes() {
		try {
			const blob = new Blob([body], {
				type: contentType || 'application/octet-stream'
			});
			const url = URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = 'request-body';
			document.body.appendChild(a);
			a.click();
			document.body.removeChild(a);
			URL.revokeObjectURL(url);
		} catch (error) {
			console.error('Error downloading bytes:', error);
		}
	}
</script>

<p>Size: {bodyByteLength} bytes</p>
<pre>{hexPreview}...</pre>
<button onclick={downloadBytes} disabled={body === ''}>Download as file</button>
