<script lang="ts">
	import Alert from '$lib/components/ui/Alert.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import KayValueList from '$lib/components/ui/KeyValueList.svelte';
	import { parseSourceDTO } from '$lib/dto-parsers';
	import type { SourceDTO } from '$lib/types';
	import { untrack } from 'svelte';
	import ListEvents from './ListEvents.svelte';
	import TestWebhook from './TestWebhook.svelte';

	type Props = {
		/** Source id shown on the details page. */
		sourceID: string;
	};

	/** Filters used to load source details. */
	type SourceFilters = {
		/** Source id shown on the details page. */
		sourceID: string;
	};

	let { sourceID }: Props = $props();

	/** Loaded source details. */
	let data = $state<SourceDTO | null>(null);

	/** Tracks whether the source details are loading. */
	let loading = $state(false);

	/** Holds the latest source loading error. */
	let error = $state<string | null>(null);

	/**
	 * Collects the current filters into a single object for easier passing to fetch functions.
	 *
	 * @returns Current source filters.
	 */
	function getCurrentFilters(): SourceFilters {
		return { sourceID };
	}

	/**
	 * Fetches the current source details from the API.
	 *
	 * @param filters - Filters to apply when fetching the source details.
	 */
	async function fetchSource(filters: SourceFilters): Promise<void> {
		loading = true;
		error = null;
		try {
			const response = await fetch(`/api/sources/${filters.sourceID}`);
			if (!response.ok) {
				throw new Error(`Failed to fetch source: ${response.statusText}`);
			}
			const rawData = await response.json();
			data = parseSourceDTO(rawData);
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error fetching source:', err);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		const filters = getCurrentFilters();

		untrack(() => {
			void fetchSource(filters);
		});
	});
</script>

{#if loading}
	<Alert>Loading source data...</Alert>
{:else if error}
	<Alert variant="error" title="Error" class="bg-surface">{error}</Alert>
{:else if data}
	<div class="flex flex-col gap-8">
		<section class="rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8">
			<div class="flex flex-col gap-6">
				<div class="flex flex-col gap-3">
					<div class="flex flex-wrap items-center gap-3">
						<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">
							Configuration
						</p>
						<Badge variant="neutral" appearance="soft" class="bg-elevated">{data.status}</Badge>
					</div>
					<h2 class="text-3xl font-semibold tracking-tight text-fg">{data.id}</h2>
					<p class="max-w-3xl text-sm leading-6 text-muted sm:text-base">
						{data.description || 'No description provided for this source.'}
					</p>
				</div>

				<div class="grid gap-4 sm:grid-cols-2">
					<div class="rounded-md border border-border-muted bg-elevated p-4">
						<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">Ingress URL</p>
						<p class="mt-2 break-all text-sm leading-6 text-fg">{data.ingress_url}</p>
					</div>
					<div class="rounded-md border border-border-muted bg-elevated p-4">
						<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">Egress URL</p>
						<p class="mt-2 break-all text-sm leading-6 text-fg">{data.egress_url}</p>
					</div>
				</div>

				<div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,0.9fr)]">
					<div class="rounded-md border border-border-muted bg-elevated p-4">
						<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">
							Static headers
						</p>
						{#if Object.keys(data.static_headers ?? {}).length > 0}
							<div class="mt-3 flex flex-col gap-2">
								{#each Object.entries(data.static_headers ?? {}) as [key, value] (key)}
									<div
										class="flex flex-col gap-1 rounded-md border border-border-muted bg-surface px-3 py-2 sm:flex-row sm:items-start sm:justify-between sm:gap-4"
									>
										<span class="text-sm font-medium text-fg">{key}</span>
										<span class="break-all text-sm text-muted">{value}</span>
									</div>
								{/each}
							</div>
						{:else}
							<p class="mt-2 text-sm text-muted">No static headers configured.</p>
						{/if}
					</div>

					<div class="rounded-md border border-border-muted bg-elevated p-4">
						<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">Metadata</p>
						<KayValueList
							items={[
								{ label: 'Status reason', value: data.status_reason },
								{ label: 'Created at', value: new Date(data.created_at).toLocaleString() },
								{ label: 'Updated at', value: new Date(data.updated_at).toLocaleString() },
								{
									label: 'Disabled at',
									value: data.disable_at ? new Date(data.disable_at).toLocaleString() : 'N/A'
								}
							]}
						/>
					</div>
				</div>
			</div>
		</section>

		<section class="rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8">
			<div class="mb-6">
				<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">Testing</p>
				<h3 class="mt-4 text-2xl font-semibold tracking-tight text-fg">Send a test webhook</h3>
				<p class="mt-3 max-w-2xl text-sm leading-6 text-muted sm:text-base">
					Exercise the ingest endpoint directly from the UI and verify the event appears below.
				</p>
			</div>
			<TestWebhook publicID={data.public_id} staticHeaders={data.static_headers} />
		</section>

		<ListEvents {sourceID} />
	</div>
{/if}
