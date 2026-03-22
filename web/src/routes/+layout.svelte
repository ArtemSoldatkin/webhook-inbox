<script lang="ts">
	import { resolve } from '$app/paths';
	import favicon from '$lib/assets/favicon.svg';
	import '$lib/highlight-init';
	import type { Snippet } from 'svelte';

	type Props = {
		/** Route content rendered inside the shared layout. */
		children: Snippet<[]>;
	};

	let { children }: Props = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<header>
	<h1>Webhook Inbox</h1>
	<nav>
		<a href={resolve('/')}>Home</a>
		<a href={resolve('/sources')}>Sources</a>
	</nav>
</header>
<main>
	{@render children()}
</main>
<footer>Webhook Inbox is an open-source project licensed under the MIT License.</footer>
