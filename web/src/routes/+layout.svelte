<script lang="ts">
	import { resolve } from '$app/paths';
	import favicon from '$lib/assets/favicon.svg';
	import '$lib/highlight-init';
	import type { Snippet } from 'svelte';
	import '../app.css';
	import ThemeToggle from './components/ThemeToggle.svelte';

	type Props = {
		/** Route content rendered inside the shared layout. */
		children: Snippet<[]>;
	};

	let { children }: Props = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div>
	<header>
		<div>
			<div>
				<a href={resolve('/')}>Webhook Inbox</a>
				<p>Inspect webhook traffic without guessing.</p>
			</div>
			<nav>
				<a href={resolve('/')}>Home</a>
				<a href={resolve('/sources')}>Sources</a>
			</nav>
			<ThemeToggle />
		</div>
	</header>
	<main>
		{@render children()}
	</main>
	<footer>
		<div>Webhook Inbox is an open-source project licensed under the MIT License.</div>
	</footer>
</div>
