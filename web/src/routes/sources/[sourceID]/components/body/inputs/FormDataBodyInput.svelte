<script lang="ts">
	import FormConstructor from './form-constructor/FormConstructor.svelte';
	import type { FormField } from './types';

	type Props = {
		/** Bound multipart payload built from form fields. */
		body: FormData;

		/** Validation error shown by the input. */
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

	/** Dynamic fields used to build the multipart payload. */
	let fields = $state<FormField[]>([]);

	$effect(() => {
		const formData = new FormData();

		for (const field of fields) {
			if (!field.name) continue;

			if (field.type === 'file' && field.value instanceof FileList) {
				for (const file of field.value) {
					formData.append(field.name, file);
				}
			} else if (field.type === 'checkbox') {
				formData.append(field.name, field.value ? 'on' : 'off');
			} else {
				formData.append(field.name, String(field.value));
			}
		}

		body = formData;
		error = null;
	});
</script>

<FormConstructor bind:fields />
