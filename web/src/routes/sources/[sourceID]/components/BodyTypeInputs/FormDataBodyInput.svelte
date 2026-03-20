<script lang="ts">
	import FormConstructor from '../FormConstructor/FormConstructor.svelte';
	import type { FormField } from '../types';

	type Props = {
		body: FormData;
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();

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
