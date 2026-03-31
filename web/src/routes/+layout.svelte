<script lang="ts">
	import { resolve } from '$app/paths';
	import favicon from '$lib/assets/favicon.svg';
	import Link from '$lib/components/ui/Link.svelte';
	import '$lib/highlight-init';
	import type { Snippet } from 'svelte';
	import '../app.css';
	import ThemeToggle from './components/ThemeToggle.svelte';

	type Props = {
		/** Route content rendered inside the shared layout. */
		children: Snippet;
	};

	let { children }: Props = $props();
	let mainContentElement = $state<HTMLElement | null>(null);

	function skipToMainContent(): void {
		mainContentElement?.focus();
		mainContentElement?.scrollIntoView({ block: 'start' });
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="min-h-screen bg-bg text-fg">
	<button
		type="button"
		onclick={skipToMainContent}
		class="sr-only absolute left-4 top-4 z-50 rounded-md bg-primary px-4 py-2 text-sm font-medium text-inverted shadow-sm focus:not-sr-only focus:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
	>
		Skip to main content
	</button>
	<header class="border-b border-border-muted bg-surface/95 shadow-sm backdrop-blur">
		<div
			class="mx-auto flex max-w-7xl items-center justify-between gap-6 px-6 py-4 sm:px-8 lg:px-12"
		>
			<div class="flex min-w-0 flex-1 flex-col gap-1">
				<Link href={resolve('/')} variant="brand">Webhook Inbox</Link>
				<p class="text-sm text-muted">Inspect webhook traffic without guessing.</p>
			</div>
			<nav
				class="flex items-center gap-1 rounded-lg border border-border bg-elevated p-1 shadow-sm"
				aria-label="Primary"
			>
				<Link href={resolve('/')} variant="nav">Home</Link>
				<Link href={resolve('/sources')} variant="nav">Sources</Link>
			</nav>
			<div class="shrink-0">
				<ThemeToggle />
			</div>
		</div>
	</header>
	<main
		bind:this={mainContentElement}
		tabindex="-1"
		class="mx-auto w-full max-w-7xl flex-1 px-6 py-8 sm:px-8 lg:px-12 lg:py-10"
	>
		{@render children()}
	</main>
	<footer class="border-t border-border-muted bg-surface">
		<div class="mx-auto max-w-7xl px-6 py-6 text-sm text-muted sm:px-8 lg:px-12">
			Webhook Inbox is an open-source project licensed under the MIT License.
		</div>
	</footer>
</div>
