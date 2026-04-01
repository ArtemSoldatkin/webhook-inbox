<script lang="ts">
	let {
		fields = $bindable()
	}: {
		fields: Array<Record<string, unknown>>;
	} = $props();

	function loadTextField(): void {
		fields = [
			{ type: 'text', name: 'title', value: 'hello' },
			{ type: 'number', name: 'count', value: 3 }
		];
	}

	function loadCheckboxField(): void {
		fields = [{ type: 'checkbox', name: 'enabled', value: true }];
	}

	function loadFileField(): void {
		const file = new File(['data'], 'demo.txt', { type: 'text/plain' });
		const fileList = Object.setPrototypeOf(
			{
				0: file,
				length: 1,
				item: (index: number) => (index === 0 ? file : null),
				[Symbol.iterator]: function* () {
					yield file;
				}
			},
			FileList.prototype
		) as FileList;

		fields = [{ type: 'file', name: 'attachment', value: fileList }];
	}
</script>

<button type="button" aria-label="load-text-field" onclick={loadTextField}>load-text-field</button>
<button type="button" aria-label="load-checkbox-field" onclick={loadCheckboxField}>
	load-checkbox-field
</button>
<button type="button" aria-label="load-file-field" onclick={loadFileField}>load-file-field</button>
