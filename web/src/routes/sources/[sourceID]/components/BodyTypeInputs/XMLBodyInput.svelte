<script lang="ts">
	type Props = {
		body: string;
		error: string | null;
	};

	let { body = $bindable(), error = $bindable() }: Props = $props();
</script>

<p>Body XML Input</p>
