<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import type { ContentType } from '$lib/types';

	type Props = {
		/** Binary body content to render. */
		body: string;

		/** Content type used when downloading the body. */
		contentType?: ContentType;
	};

	let { body, contentType }: Props = $props();

	/** Size of the current binary body in bytes. */
	const bodyByteLength = $derived(body.length);

	/** Hex preview of the first bytes in the body. */
	const hexPreview = $derived(
		Array.from(body.slice(0, 16))
			.map((c) => c.charCodeAt(0).toString(16).padStart(2, '0'))
			.join(' ')
	);

	/**
	 * Converts a binary string into a byte array.
	 *
	 * @param str - Binary string to convert.
	 * @returns Byte array for the provided string.
	 */
	function binaryStringToUint8Array(str: string): Uint8Array<ArrayBuffer> {
		const buffer = new ArrayBuffer(str.length);
		const bytes = new Uint8Array(buffer);

		for (let i = 0; i < str.length; i++) {
			bytes[i] = str.charCodeAt(i) & 0xff;
		}

		return bytes;
	}

	/** Downloads the current binary body as a file. */
	function downloadBytes(): void {
		try {
			const bytes = binaryStringToUint8Array(body);
			const blob = new Blob([bytes], {
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
<Button onclick={downloadBytes} disabled={body === ''}>Download as file</Button>
