<script lang="ts">
	import { resolve } from '$app/paths';
	import Link from '$lib/components/ui/Link.svelte';
	import env from '$lib/env';

	/**
	 * The quick start steps should match the actual flow in the app,
	 * and the example request should be easy to copy and paste when they are ready to send test traffic.
	 */
	const quickStartSteps = [
		{
			label: '01',
			title: 'Create a source',
			description:
				'Define an egress URL and optional static headers for the integration you want to observe.'
		},
		{
			label: '02',
			title: 'Send a webhook',
			description:
				'Use the generated ingest endpoint with curl, a provider like Stripe, or the built-in tester.'
		},
		{
			label: '03',
			title: 'Inspect the result',
			description:
				'Review stored payloads, headers, timestamps, and delivery attempts from the source detail view.'
		}
	] as const;

	/**
	 * The captured items should reflect the actual data we store for each event,
	 * and the features should reflect the actual capabilities of the app.
	 * Avoid aspirational language that oversells or misrepresents what we currently offer.
	 */
	const capturedItems = [
		'Request body',
		'Headers and query params',
		'Event timestamps',
		'Delivery attempts',
		'Source-specific metadata'
	] as const;

	const features = [
		{
			title: 'Source-based endpoints',
			description: 'Keep integrations isolated by assigning each workflow its own ingress URL.'
		},
		{
			title: 'Full payload capture',
			description: 'Store request bodies and headers so you can inspect exactly what arrived.'
		},
		{
			title: 'Cursor pagination',
			description: 'Browse large event streams without losing ordering or context.'
		},
		{
			title: 'Delivery tracking',
			description: 'Follow retry behavior and delivery outcomes for each recorded event.'
		},
		{
			title: 'Built-in testing',
			description: 'Send sample webhook traffic from the UI before you point a real provider at it.'
		},
		{
			title: 'Debug-first workflow',
			description: 'Move from ingestion to inspection fast, without building ad hoc logging first.'
		}
	] as const;

	/**
	 * The example request should be a realistic sample that users can easily modify with their own source key and event data.
	 * It should demonstrate how to send a POST request with a JSON body to the ingest endpoint, which is the most common way webhooks are sent.
	 * Avoid including unnecessary flags or complex payloads that might overwhelm users who are just trying to get started.
	 */
	const exampleRequest = `curl -X POST ${env.VITE_API_BASE_URL}/api/v1/ingest/YOUR_SOURCE_KEY \\
  -H "Content-Type: application/json" \\
  -d '{"event":"payment.succeeded","source":"demo"}'`;
</script>

<svelte:head>
	<title>Webhook Inbox</title>
	<meta
		name="description"
		content="Create unique webhook endpoints, capture incoming requests, and inspect deliveries in one place."
	/>
</svelte:head>

<div class="min-h-screen bg-bg text-fg">
	<div class="mx-auto flex max-w-7xl flex-col gap-20 px-6 py-12 sm:px-8 lg:px-12 lg:py-16">
		<section class="grid gap-8 lg:grid-cols-[minmax(0,1.2fr)_minmax(22rem,0.8fr)] lg:items-start">
			<div class="flex flex-col gap-8">
				<div
					class="inline-flex w-fit items-center rounded-md border border-border bg-elevated px-3 py-1 text-sm font-medium text-muted shadow-sm"
				>
					Minimal webhook recorder
				</div>

				<div class="flex max-w-3xl flex-col gap-5">
					<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">
						Receive, inspect, and replay webhook traffic
					</p>
					<h1 class="max-w-3xl text-4xl font-semibold tracking-tight text-fg sm:text-5xl">
						See every webhook the way your provider actually sent it.
					</h1>
					<p class="max-w-2xl text-base leading-7 text-muted sm:text-lg">
						Create unique endpoints, capture request bodies and headers, and debug delivery flows
						without adding temporary logging to your services.
					</p>
				</div>

				<div class="flex flex-col gap-3 sm:flex-row">
					<Link href={resolve('/sources')} variant="primary">
						Open Sources
					</Link>
					<Link href="#how-it-works" variant="secondary">
						How It Works
					</Link>
				</div>

				<div class="grid gap-4 sm:grid-cols-3">
					<div class="rounded-lg border border-border bg-surface p-5 shadow-sm">
						<p class="text-2xl font-semibold tracking-tight text-fg">202</p>
						<p class="mt-2 text-sm leading-6 text-muted">
							Fast ingest responses keep senders moving.
						</p>
					</div>
					<div class="rounded-lg border border-border bg-surface p-5 shadow-sm">
						<p class="text-2xl font-semibold tracking-tight text-fg">JSONB</p>
						<p class="mt-2 text-sm leading-6 text-muted">
							Payloads and headers stay queryable in Postgres.
						</p>
					</div>
					<div class="rounded-lg border border-border bg-surface p-5 shadow-sm">
						<p class="text-2xl font-semibold tracking-tight text-fg">UI + API</p>
						<p class="mt-2 text-sm leading-6 text-muted">
							Create sources, inspect events, and send test traffic.
						</p>
					</div>
				</div>
			</div>

			<aside class="rounded-lg border border-border bg-surface p-6 shadow-md sm:p-7">
				<div class="flex flex-col gap-6">
					<div class="flex items-center justify-between gap-4">
						<div>
							<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">
								Example request
							</p>
							<p class="mt-2 text-sm leading-6 text-muted">
								Point a provider or a local script at your generated ingest URL and inspect the
								event a few seconds later.
							</p>
						</div>
						<div
							class="hidden rounded-md border border-border bg-elevated px-3 py-1 text-xs font-medium text-subtle shadow-sm sm:block"
						>
							curl
						</div>
					</div>

					<pre
						class="overflow-x-auto rounded-md border border-border-muted bg-elevated p-4 text-sm leading-6 text-fg shadow-sm"><code
							>{exampleRequest}</code
						></pre>

					<div class="rounded-md border border-border-muted bg-elevated p-4">
						<p class="text-sm font-medium text-fg">Typical use cases</p>
						<p class="mt-2 text-sm leading-6 text-muted">
							Stripe callbacks, GitHub webhooks, Twilio events, internal service integrations, and
							local contract testing.
						</p>
					</div>
				</div>
			</aside>
		</section>

		<section id="how-it-works" class="border-t border-border-muted pt-16">
			<div class="flex flex-col gap-4 lg:max-w-3xl">
				<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">Start in 3 steps</p>
				<h2 class="text-3xl font-semibold tracking-tight text-fg sm:text-4xl">
					A homepage that acts like onboarding.
				</h2>
				<p class="text-base leading-7 text-muted">
					The core flow already exists in the app. The landing page should simply point users into
					it without forcing them to read the README first.
				</p>
			</div>

			<div class="mt-8 grid gap-4 lg:grid-cols-3">
				{#each quickStartSteps as step (step.label)}
					<article class="rounded-lg border border-border bg-surface p-6 shadow-sm">
						<p class="text-sm font-semibold tracking-[0.18em] text-primary">{step.label}</p>
						<h3 class="mt-4 text-xl font-semibold tracking-tight text-fg">{step.title}</h3>
						<p class="mt-3 text-sm leading-6 text-muted">{step.description}</p>
					</article>
				{/each}
			</div>
		</section>

		<section class="grid gap-6 lg:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)]">
			<article class="rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8">
				<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">
					What gets captured
				</p>
				<h2 class="mt-4 text-3xl font-semibold tracking-tight text-fg">
					Everything you need to debug the handoff.
				</h2>
				<p class="mt-4 text-base leading-7 text-muted">
					Webhook Inbox is most useful when something subtle breaks. The recorded event should let
					you compare assumptions against the real request.
				</p>
				<ul class="mt-8 space-y-3">
					{#each capturedItems as item (item)}
						<li
							class="flex items-start gap-3 rounded-md border border-border-muted bg-elevated px-4 py-3"
						>
							<span class="mt-1 h-2.5 w-2.5 shrink-0 rounded-full bg-primary"></span>
							<span class="text-sm leading-6 text-fg">{item}</span>
						</li>
					{/each}
				</ul>
			</article>

			<article class="rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8">
				<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">
					Built for inspection and debugging
				</p>
				<div class="mt-6 grid gap-4 sm:grid-cols-2">
					{#each features as feature (feature.title)}
						<section class="rounded-md border border-border-muted bg-elevated p-5">
							<h3 class="text-lg font-semibold tracking-tight text-fg">{feature.title}</h3>
							<p class="mt-3 text-sm leading-6 text-muted">{feature.description}</p>
						</section>
					{/each}
				</div>
			</article>
		</section>

		<section>
			<div class="rounded-lg bg-primary px-6 py-8 text-inverted shadow-lg sm:px-8 sm:py-10">
				<div class="flex flex-col gap-8 lg:flex-row lg:items-end lg:justify-between">
					<div class="max-w-3xl">
						<p class="text-sm font-medium uppercase tracking-[0.18em] text-inverted">
							Ready to start?
						</p>
						<h2 class="mt-4 text-3xl font-semibold tracking-tight text-inverted sm:text-4xl">
							Create a source, send a sample payload, and inspect the first event.
						</h2>
						<p class="mt-4 max-w-2xl text-base leading-7 text-inverted/80">
							The sources workflow is already the operational center of the app. Use the homepage to
							funnel people there with enough context to know what they are looking at.
						</p>
					</div>
					<div class="flex flex-col gap-3 sm:flex-row">
						<Link href={resolve('/sources')} variant="secondary">
							Go To Sources
						</Link>
						<Link href={resolve('/sources')} variant="inverse">
							Create Your First Source
						</Link>
					</div>
				</div>
			</div>
		</section>
	</div>
</div>
