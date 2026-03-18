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
	}
</script>

<p>Size: {bodyByteLength} bytes</p>
<pre>{hexPreview}...</pre>
<button on:click={downloadBytes} disabled={body === ''}>Download as file</button>
