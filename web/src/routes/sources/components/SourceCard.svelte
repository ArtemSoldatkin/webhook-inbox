<script lang="ts">
	import { resolve } from '$app/paths';
	import Alert from '$lib/components/ui/Alert.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Eyebrow from '$lib/components/ui/Eyebrow.svelte';
	import KeyValueList from '$lib/components/ui/KeyValueList.svelte';
	import Link from '$lib/components/ui/Link.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import type { SourceDTO } from '$lib/types';
	import Icon from '@iconify/svelte';

	type Props = {
		/** Source data to display within this card. */
		source: SourceDTO;

		/** Whether the source ID should be rendered as a link to the source details page. Defaults to false. */
		idAsLink?: boolean;

		/** Optional callback to trigger a refresh of the sources list after updating source status. */
		onStatusUpdate?: () => void;
	};

	let { source, idAsLink = false, onStatusUpdate }: Props = $props();

	/** Local state for managing edit mode, form input, loading state, and error messages. */
	let isEditing = $state(false);
	/** Local state for the new status value when editing. */
	let newStatus = $state<string>(source.status);

	/** Local state for tracking loading status during async operations. */
	let loading = $state(false);
	/** Local state for storing error messages to display in the UI. */
	let error = $state<string | null>(null);

	/**
	 * Async function to handle updating the source status via API call.
	 *
	 * @param newStatus - The new status to set for the source.
	 */
	async function updateSourceStatus(newStatus: string): Promise<void> {
		loading = true;
		error = null;
		try {
			const response = await fetch(`/api/sources/${source.id}/status`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					status: newStatus,
					status_reason: `Status changed to ${newStatus} via UI`
				})
			});
			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.message || 'Failed to update source status');
			}
			source.status = newStatus;
			isEditing = false;
			onStatusUpdate?.();
		} catch (err) {
			console.error('Failed to update source status:', err);
			error = 'Failed to update source status. Please try again.';
			newStatus = source.status;
		} finally {
			loading = false;
		}
	}

	/**
	 * Toggles the edit mode for the source status. If currently editing and the status has changed, it triggers an update.
	 * Otherwise, it simply toggles the edit mode and resets the newStatus to the current source status.
	 */
	function toggleEditing(): void {
		if (isEditing && newStatus !== source.status) {
			updateSourceStatus(newStatus);
		} else {
			newStatus = source.status;
			isEditing = !isEditing;
		}
	}
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
			{#if isEditing}
				<div class="flex flex-col lg:flex-row items-start lg:items-center gap-1">
					<Select
						bind:value={newStatus}
						options={[
							{ value: 'active', label: 'Active' },
							{ value: 'paused', label: 'Paused' },
							{ value: 'quarantined', label: 'Quarantined' },
							{ value: 'disabled', label: 'Disabled' }
						]}
					/>
					{#if error}
						<Alert variant="error">
							{error}
						</Alert>
					{/if}
				</div>
			{:else}
				<Badge variant="neutral" appearance="soft">{source.status}</Badge>
			{/if}
			<Button
				variant="ghost"
				class="ml-auto"
				onclick={toggleEditing}
				disabled={loading}
				aria-label={isEditing ? 'Save source status' : 'Edit source status'}
			>
				{#if isEditing}
					<Icon
						icon="material-symbols:save-as-outline-rounded"
						class="text-xl"
						aria-hidden="true"
					/>
				{:else}
					<Icon
						icon="material-symbols:edit-document-outline-rounded"
						class="text-xl"
						aria-hidden="true"
					/>
				{/if}
			</Button>
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
				<div
					class="mt-3 flex flex-col divide-y divide-border-muted rounded-md border border-border-muted"
				>
					{#each Object.entries(source.static_headers ?? {}) as [key, value] (key)}
						<div
							class="flex flex-col gap-1 px-3 py-3 sm:flex-row sm:items-start sm:justify-between sm:gap-4"
						>
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
