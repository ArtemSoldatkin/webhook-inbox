<script lang="ts">
	import Alert from '$lib/components/ui/Alert.svelte';
	import SectionHeader from '$lib/components/ui/SectionHeader.svelte';
	import { parseSourceDTO } from '$lib/dto-parsers';
	import type { SourceDTO } from '$lib/types';
	import { untrack } from 'svelte';
	import SourceCard from '../../components/SourceCard.svelte';
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
			<SourceCard source={data} />
		</section>

		<section class="rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8">
			<SectionHeader
				eyebrow="Testing"
				title="Send a test webhook"
				description="Exercise the ingest endpoint directly from the UI and verify the event appears below."
				titleAs="h3"
				class="mb-6"
			/>
			<TestWebhook publicID={data.public_id} staticHeaders={data.static_headers} />
		</section>

		<ListEvents {sourceID} />
	</div>
{:else}
	<Alert variant="warning" title="Source not found" class="bg-surface">
		The requested source could not be found. It may have been deleted.
	</Alert>
{/if}
