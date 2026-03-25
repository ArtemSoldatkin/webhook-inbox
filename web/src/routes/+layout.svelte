<script lang="ts">
	import { resolve } from '$app/paths';
	import Link from '$lib/components/ui/Link.svelte';
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

<div class="min-h-screen bg-bg text-fg">
	<header class="border-b border-border-muted bg-surface/95 shadow-sm backdrop-blur">
		<div class="mx-auto flex max-w-7xl items-center justify-between gap-6 px-6 py-4 sm:px-8 lg:px-12">
			<div class="flex min-w-0 flex-1 flex-col gap-1">
				<Link href={resolve('/')} variant="brand">Webhook Inbox</Link>
				<p class="text-sm text-muted">Inspect webhook traffic without guessing.</p>
			</div>
			<nav class="flex items-center gap-1 rounded-lg border border-border bg-elevated p-1 shadow-sm">
				<Link href={resolve('/')} variant="nav">Home</Link>
				<Link href={resolve('/sources')} variant="nav">Sources</Link>
			</nav>
			<div class="shrink-0">
				<ThemeToggle />
			</div>
		</div>
	</header>
	<main class="mx-auto w-full max-w-7xl flex-1 px-6 py-8 sm:px-8 lg:px-12 lg:py-10">
		{@render children()}
	</main>
	<footer class="border-t border-border-muted bg-surface">
		<div class="mx-auto max-w-7xl px-6 py-6 text-sm text-muted sm:px-8 lg:px-12">
			Webhook Inbox is an open-source project licensed under the MIT License.
		</div>
	</footer>
</div>
