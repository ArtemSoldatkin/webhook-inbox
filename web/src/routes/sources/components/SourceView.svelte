<script lang="ts">
	import { resolve } from '$app/paths';
	import KayValueList from '$lib/components/ui/KayValueList.svelte';
	import Link from '$lib/components/ui/Link.svelte';
	import { type SourceDTO } from '$lib/types';

	type Props = {
		/** Source row rendered in the sources list. */
		source: SourceDTO;
	};

	let { source }: Props = $props();
</script>

<section class="rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8">
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-3">
			<div class="flex flex-wrap items-center gap-3">
				<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">Source</p>
				<span
					class="inline-flex w-fit rounded-full border border-border bg-elevated px-3 py-1 text-xs font-medium uppercase tracking-[0.12em] text-muted"
				>
					{source.status}
				</span>
			</div>
			<h1 class="text-3xl font-semibold tracking-tight text-fg">
				<Link href={resolve(`/sources/${source.id}`)} variant="inline">{source.id}</Link>
			</h1>
			<p class="max-w-3xl text-sm leading-6 text-muted sm:text-base">
				{source.description || 'No description provided for this source.'}
			</p>
		</div>

		<div class="grid gap-4 sm:grid-cols-2">
			<div class="rounded-md border border-border-muted bg-elevated p-4">
				<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">Ingress URL</p>
				<p class="mt-2 break-all text-sm leading-6 text-fg">{source.ingress_url}</p>
			</div>
			<div class="rounded-md border border-border-muted bg-elevated p-4">
				<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">Egress URL</p>
				<p class="mt-2 break-all text-sm leading-6 text-fg">{source.egress_url}</p>
			</div>
		</div>

		<div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,0.9fr)]">
			<div class="rounded-md border border-border-muted bg-elevated p-4">
				<p class="text-xs font-medium uppercase tracking-[0.12em] text-subtle">Static headers</p>
				{#if Object.keys(source.static_headers ?? {}).length > 0}
					<div class="mt-3 flex flex-col gap-2">
						{#each Object.entries(source.static_headers ?? {}) as [key, value] (key)}
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
						{ label: 'Status reason', value: source.status_reason },
						{ label: 'Created at', value: new Date(source.created_at).toLocaleString() },
						{ label: 'Updated at', value: new Date(source.updated_at).toLocaleString() },
						{
							label: 'Disabled at',
							value: source.disable_at ? new Date(source.disable_at).toLocaleString() : 'N/A'
						}
					]}
				/>
			</div>
		</div>
	</div>
</section>
