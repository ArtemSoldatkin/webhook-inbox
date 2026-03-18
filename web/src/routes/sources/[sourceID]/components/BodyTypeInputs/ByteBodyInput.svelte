<script lang="ts">
	type Props = {
		body: string;
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

	function handleFileChange(event: Event) {
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

	function isValidBase64(str: string) {
		try {
			return btoa(atob(str)) === str.replace(/\s/g, '');
		} catch {
			return false;
		}
	}

	$effect(() => {
		if (body && !isValidBase64(body)) {
			error = 'Invalid base64 string';
		} else {
			error = null;
		}
	});
</script>

<textarea rows={5} bind:value={body} placeholder="Enter base64 body here..."></textarea>
<input type="file" onchange={(e) => handleFileChange(e)} />
