<script lang="ts">
	let {
		body = $bindable(),
		error = $bindable()
	}: { body?: string | FormData; error?: string | null } = $props();
</script>

<div data-testid="body-input-child">
	<span>{typeof body === 'string' ? body : 'form-data'}</span>
	<span>{error ?? ''}</span>
</div>
