<script lang="ts">
	let { body = '', contentType = '' }: { body?: string; contentType?: string } = $props();
</script>

<div data-testid="body-view-child">
	<span>{body}</span>
	<span>{contentType}</span>
</div>
