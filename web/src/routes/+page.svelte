<script lang="ts">
	import { resolve } from '$app/paths';
	import env from '$lib/env';

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

<section>
	<div>
		<div>Minimal webhook recorder</div>
		<div>
			<p>Receive, inspect, and replay webhook traffic</p>
			<h1>See every webhook the way your provider actually sent it.</h1>
			<p>
				Create unique endpoints, capture request bodies and headers, and debug delivery flows
				without adding temporary logging to your services.
			</p>
		</div>
		<div>
			<a href={resolve('/sources')}>Open Sources</a>
			<a href="#how-it-works">How It Works</a>
		</div>
		<div>
			<div>
				<p>202</p>
				<p>Fast ingest responses keep senders moving.</p>
			</div>
			<div>
				<p>JSONB</p>
				<p>Payloads and headers stay queryable in Postgres.</p>
			</div>
			<div>
				<p>UI + API</p>
				<p>Create sources, inspect events, and send test traffic.</p>
			</div>
		</div>
	</div>

	<aside>
		<div>Example request</div>
		<div>
			<p>
				Point a provider or a local script at your generated ingest URL and inspect the event a few
				seconds later.
			</p>
			<pre><code>{exampleRequest}</code></pre>
			<div>
				<p>Typical use cases</p>
				<p>
					Stripe callbacks, GitHub webhooks, Twilio events, internal service integrations, and local
					contract testing.
				</p>
			</div>
		</div>
	</aside>
</section>

<section id="how-it-works">
	<div>
		<div>
			<p>Start in 3 steps</p>
			<h2>A homepage that acts like onboarding.</h2>
		</div>
		<p>
			The core flow already exists in the app. The landing page should simply point users into it
			without forcing them to read the README first.
		</p>
	</div>
	<div>
		{#each quickStartSteps as step (step.label)}
			<article>
				<p>{step.label}</p>
				<h3>{step.title}</h3>
				<p>{step.description}</p>
			</article>
		{/each}
	</div>
</section>

<section>
	<article>
		<p>What gets captured</p>
		<h2>Everything you need to debug the handoff.</h2>
		<p>
			Webhook Inbox is most useful when something subtle breaks. The recorded event should let you
			compare assumptions against the real request.
		</p>
		<ul>
			{#each capturedItems as item (item)}
				<li>{item}</li>
			{/each}
		</ul>
	</article>

	<article>
		<p>Built for inspection and debugging</p>
		<div>
			{#each features as feature (feature.title)}
				<section>
					<h3>{feature.title}</h3>
					<p>{feature.description}</p>
				</section>
			{/each}
		</div>
	</article>
</section>

<section>
	<div>
		<div>
			<p>Ready to start?</p>
			<h2>Create a source, send a sample payload, and inspect the first event.</h2>
			<p>
				The sources workflow is already the operational center of the app. Use the homepage to
				funnel people there with enough context to know what they are looking at.
			</p>
		</div>
		<div>
			<a href={resolve('/sources')}>Go To Sources</a>
			<a href={resolve('/sources')}>Create Your First Source</a>
		</div>
	</div>
</section>
