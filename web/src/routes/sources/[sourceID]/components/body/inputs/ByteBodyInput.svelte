<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import FileInput from '$lib/components/ui/FileInput.svelte';
	import Textarea from '$lib/components/ui/Textarea.svelte';

	type Props = {
		/** Bound base64-encoded body value. */
		body: string;

		/** Validation error shown by the input. */
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

	/**
	 * Reads a selected file and stores it as base64 text.
	 *
	 * @param event - File input change event.
	 */
	function handleFileChange(event: Event): void {
		const input = event.target as HTMLInputElement;
		if (input.files && input.files[0]) {
			const file = input.files[0];
			const reader = new FileReader();
			reader.onload = () => {
				const arrayBuffer = reader.result as ArrayBuffer;
				const uint8Array = new Uint8Array(arrayBuffer);
				const binary = Array.from(uint8Array)
					.map((b) => String.fromCharCode(b))
					.join('');
				body = btoa(binary);
				error = null;
			};
			reader.onerror = () => {
				error = 'Failed to read file';
			};
			reader.readAsArrayBuffer(file);
		}
	}

	/**
	 * Checks whether a string is valid base64.
	 *
	 * @param str - String to validate.
	 * @returns Whether the string can be decoded as base64.
	 */
	function isValidBase64(str: string): boolean {
		try {
			return btoa(atob(str)) === str.replace(/\s/g, '');
		} catch {
			return false;
		}
	}

	/** Clears the current byte body input. */
	function handleClear(): void {
		body = '';
		error = null;
	}

	$effect(() => {
		if (body && !isValidBase64(body)) {
			error = 'Invalid base64 string';
		} else {
			error = null;
		}
	});
</script>

<div class="flex flex-col gap-4">
	<Textarea
		aria-label="Base64 request body"
		rows={10}
		bind:value={body}
		placeholder="Enter base64 body here..."
	/>
	<div class="rounded-md border border-border-muted bg-elevated p-4">
		<label class="text-sm font-medium text-fg">
			Upload file
			<FileInput
				onchange={(e) => handleFileChange(e)}
				class="mt-2 block w-full text-sm"
				aria-label="Upload file"
			/>
		</label>
	</div>
	<div class="flex justify-end">
		<Button type="button" onclick={handleClear} disabled={!body} variant="secondary">Clear</Button>
	</div>
</div>
