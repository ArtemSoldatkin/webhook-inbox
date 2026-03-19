<script lang="ts">
	import formatXML from 'xml-formatter';

	type Props = {
		body: string;
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

	function validateXML(xml: string): string | null {
		if (!xml.trim()) {
			return null;
		}

		const parser = new DOMParser();
		const doc = parser.parseFromString(xml, 'application/xml');
		const parserError = doc.querySelector('parsererror');

		return parserError ? 'Invalid XML format' : null;
	}

	function formatInput() {
		try {
			body = formatXML(body);
		} catch (err) {
			console.error('Error formatting XML body:', err);
			error = 'Cannot format invalid XML';
		}
	}

	function handleClear() {
		body = '';
		error = null;
	}

	$effect(() => {
		error = validateXML(body);
	});
</script>

<textarea bind:value={body} rows="10" cols="50" placeholder="Enter XML body here..."></textarea>
<button type="button" onclick={formatInput} disabled={!!error}>Format</button>
<button type="button" onclick={handleClear} disabled={!body}>Clear</button>
