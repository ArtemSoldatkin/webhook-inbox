<script lang="ts">
	import { resolve } from '$app/paths';
	import Badge from '$lib/components/ui/Badge.svelte';
	import Eyebrow from '$lib/components/ui/Eyebrow.svelte';
	import KeyValueList from '$lib/components/ui/KeyValueList.svelte';
	import Link from '$lib/components/ui/Link.svelte';
	import type { SourceDTO } from '$lib/types';

	type Props = {
		/** Source data to display within this card. */
		source: SourceDTO;

		/** Whether the source ID should be rendered as a link to the source details page. Defaults to false. */
		idAsLink?: boolean;
	};

	let { source, idAsLink = false }: Props = $props();
</script>

<article class="rounded-lg border border-border bg-surface p-5 shadow-sm">
	<div class="grid grid-cols-1 gap-5 xl:grid-cols-2">
		<div class="xl:col-span-2 flex flex-wrap items-center gap-3">
			<h3 class="text-xl font-semibold tracking-tight text-fg">
				{#if idAsLink}
					<Link href={resolve(`/sources/${source.id}`)} variant="inline">{source.id}</Link>
				{:else}
					{source.id}
				{/if}
			</h3>
			<Badge variant="neutral" appearance="soft">{source.status}</Badge>
		</div>

		<p class="xl:col-span-2 text-sm leading-6 text-muted">
			{source.description || 'No description provided.'}
		</p>

		<section class="border-t border-border-muted pt-4">
			<Eyebrow>Ingress URL</Eyebrow>
			<p class="mt-2 break-all text-sm leading-6 text-fg">{source.ingress_url}</p>
		</section>

		<section class="border-t border-border-muted pt-4">
			<Eyebrow>Egress URL</Eyebrow>
			<p class="mt-2 break-all text-sm leading-6 text-fg">{source.egress_url}</p>
		</section>

		<section class="border-t border-border-muted pt-4">
			<Eyebrow>Static headers</Eyebrow>
			{#if Object.keys(source.static_headers ?? {}).length > 0}
				<div class="mt-3 flex flex-col divide-y divide-border-muted rounded-md border border-border-muted">
					{#each Object.entries(source.static_headers ?? {}) as [key, value] (key)}
						<div class="flex flex-col gap-1 px-3 py-3 sm:flex-row sm:items-start sm:justify-between sm:gap-4">
							<span class="text-sm font-medium text-fg">{key}</span>
							<span class="break-all text-sm text-muted">{value}</span>
						</div>
					{/each}
				</div>
			{:else}
				<p class="mt-2 text-sm text-muted">No static headers configured.</p>
			{/if}
		</section>

		<section class="border-t border-border-muted pt-4">
			<Eyebrow>Metadata</Eyebrow>
			<KeyValueList
				items={[
					{ label: 'Status reason', value: source.status_reason },
					{
						label: 'Created at',
						value: new Date(source.created_at).toLocaleString()
					},
					{
						label: 'Updated at',
						value: new Date(source.updated_at).toLocaleString()
					},
					{
						label: 'Disabled at',
						value: source.disable_at ? new Date(source.disable_at).toLocaleString() : 'N/A'
					}
				]}
			/>
		</section>
	</div>
</article>
