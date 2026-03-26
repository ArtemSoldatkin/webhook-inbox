<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Eyebrow from '$lib/components/ui/Eyebrow.svelte';
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

<div class="mt-3 flex flex-col gap-4">
	<div class="rounded-md border border-border-muted bg-elevated p-4">
		<div class="grid gap-3 sm:grid-cols-2">
			<div>
				<Eyebrow>Size</Eyebrow>
				<p class="mt-2 text-sm text-fg">{bodyByteLength} bytes</p>
			</div>
			<div>
				<Eyebrow>Preview</Eyebrow>
				<pre class="mt-2 overflow-x-auto text-sm leading-6 text-fg">{hexPreview}...</pre>
			</div>
		</div>
	</div>
	<div class="flex justify-end">
		<Button onclick={downloadBytes} disabled={body === ''} variant="secondary"
			>Download as file</Button
		>
	</div>
</div>
