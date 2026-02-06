<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<header>
	<h1>Webhook Inbox</h1>
	<nav>
		<ul>
			<li><a href="/">Home</a></li>
			<li><a href="/users">Users</a></li>
			<li><a href="/endpoints">Endpoints</a></li>
		</ul>
	</nav>
</header>
<main>
	{@render children()}
</main>
<footer>
	<p>&copy; 2026 Webhook Inbox</p>
</footer>
