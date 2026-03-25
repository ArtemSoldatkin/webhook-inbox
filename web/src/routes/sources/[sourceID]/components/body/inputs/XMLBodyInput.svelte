<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import formatXML from 'xml-formatter';

	type Props = {
		/** Bound XML body value. */
		body: string;

		/** Validation error shown by the input. */
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

	/**
	 * Validates the current XML body.
	 *
	 * @param xml - XML string to validate.
	 * @returns Validation error message or `null`.
	 */
	function validateXML(xml: string): string | null {
		if (!xml.trim()) {
			return null;
		}

		const parser = new DOMParser();
		const doc = parser.parseFromString(xml, 'application/xml');
		const parserError = doc.querySelector('parsererror');

		return parserError ? 'Invalid XML format' : null;
	}

	/** Pretty-prints the current XML body when it is valid. */
	function formatInput(): void {
		try {
			body = formatXML(body);
		} catch (err) {
			console.error('Error formatting XML body:', err);
			error = 'Cannot format invalid XML';
		}
	}

	/** Clears the current XML body. */
	function handleClear(): void {
		body = '';
		error = null;
	}

	$effect(() => {
		error = validateXML(body);
	});
</script>

<div class="flex flex-col gap-4">
	<textarea
		bind:value={body}
		rows="10"
		placeholder="Enter XML body here..."
		class="min-h-56 w-full rounded-md border border-border bg-surface px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
	></textarea>
	<div class="flex flex-wrap justify-end gap-3">
		<Button type="button" onclick={formatInput} disabled={!!error} variant="secondary">Format</Button>
		<Button type="button" onclick={handleClear} disabled={!body} variant="secondary">Clear</Button>
	</div>
</div>
