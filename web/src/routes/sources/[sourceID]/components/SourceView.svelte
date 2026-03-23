<script lang="ts">
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
			fetchSource(filters);
		});
	});
</script>

{#if loading}
	<p>Loading source data...</p>
{:else if error}
	<p class="error">Error: {error}</p>
{:else if data}
	<section>
		<h2>{data.id}</h2>
		<p>{data.description}</p>
		<p>{data.ingress_url}</p>
		<p>{data.egress_url}</p>
		<p>Static headers:</p>
		{#each Object.entries(data.static_headers ?? {}) as [key, value] (key)}
			<p>{key}: {value}</p>
		{/each}
		<p>{data.status}</p>
		<p>{data.status_reason}</p>
		<p>Created at: {new Date(data.created_at).toLocaleString()}</p>
		<p>Updated at: {new Date(data.updated_at).toLocaleString()}</p>
		<p>
			Disabled at: {data.disable_at ? new Date(data.disable_at).toLocaleString() : 'N/A'}
		</p>
	</section>
	<section>
		<h3>Test Webhook</h3>
		<TestWebhook publicID={data.public_id} staticHeaders={data.static_headers} />
		<ListEvents {sourceID} />
	</section>
{/if}
