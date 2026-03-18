<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import '$lib/highlight-init';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<header>
	<h1>Webhook Inbox</h1>
	<nav>
		<a href="/">Home</a>
		<a href="/sources">Sources</a>
	</nav>
</header>
<main>
	{@render children()}
</main>
<footer>Webhook Inbox is an open-source project licensed under the MIT License.</footer>
