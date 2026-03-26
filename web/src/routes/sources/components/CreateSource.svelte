<script lang="ts">
	import InputMap from '$lib/components/InputMap.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import env from '$lib/env';
	import type { SourceDTO } from '$lib/types';

	/** Source payload used by the create form. */
	type NewSource = Omit<
		SourceDTO,
		| 'id'
		| 'public_id'
		| 'ingress_url'
		| 'status'
		| 'status_reason'
		| 'created_at'
		| 'updated_at'
		| 'disable_at'
	>;

	/** Current form data for the new source. */
	let data = $state<NewSource>(newData());

	/** Tracks whether a source creation request is in flight. */
	let loading = $state(false);

	/** Holds the latest source creation error. */
	let error = $state<string | null>(null);

	/** Validation error for the current egress URL. */
	let egressError = $derived.by<string | null>(() => {
		if (data.egress_url.trim() !== '' && !validateEgressUrl(data.egress_url)) {
			return 'Valid Egress URL is required';
		}
		return null;
	});

	/**
	 * Creates a fresh source payload with default values.
	 *
	 * @returns Empty source form data.
	 */
	function newData(): NewSource {
		return { egress_url: '', static_headers: {}, description: '' };
	}

	/** Sends the current source form to the API. */
	async function createSource(): Promise<void> {
		loading = true;
		error = null;
		try {
			const response = await fetch('/api/sources', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(data)
			});
			if (!response.ok) {
				throw new Error(`Failed to create source: ${response.statusText}`);
			}
			const newSource = await response.json();
			data = newData(); // Reset form after successful creation
			console.log('Created new source:', newSource);
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error creating source:', err);
		} finally {
			loading = false;
		}
	}

	/**
	 * Checks whether the form is ready to submit.
	 *
	 * @returns Whether the current input is valid.
	 */
	function validateInput(): boolean {
		return validateEgressUrl(data.egress_url);
	}

	/**
	 * Handles source creation form submission.
	 *
	 * @param event - Form submission event.
	 */
	function handleSubmit(event: SubmitEvent): void {
		event.preventDefault();
		const isValid = validateInput();
		if (!isValid) {
			console.warn('Invalid input, cannot create source');
			return;
		}
		createSource();
	}

	/**
	 * Validates an egress URL against protocol and network rules.
	 *
	 * @param url - Egress URL to validate.
	 * @returns Whether the URL is allowed.
	 */
	function validateEgressUrl(url: string): boolean {
		try {
			const parsedUrl = new URL(url);
			if (
				(parsedUrl.protocol !== 'http:' && parsedUrl.protocol !== 'https:') ||
				parsedUrl.href.length > 2048
			) {
				return false;
			}
			if (env.VITE_ENV === 'dev') return true;
			return (
				/^https?:\/\//.test(parsedUrl.href) &&
				!/^https?:\/\/(localhost|127\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})|0\.0\.0\.0|\[?::1\]?)(\/|:|$)/.test(
					parsedUrl.href
				) &&
				!/^https?:\/\/\[::ffff:127\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\]/.test(parsedUrl.href) &&
				!/^https?:\/\/10\./.test(parsedUrl.href) &&
				!/^https?:\/\/192\.168\./.test(parsedUrl.href) &&
				!/^https?:\/\/172\.(1[6-9]|2[0-9]|3[0-1])\./.test(parsedUrl.href) &&
				!/^https?:\/\/169\.254\.169\.254(\/|:|$)/.test(parsedUrl.href) &&
				!/^https?:\/\/\[::ffff:0\.0\.0\.0\]/.test(parsedUrl.href) &&
				!/^https?:\/\/localhost\.(\/|:|$)/.test(parsedUrl.href)
			);
		} catch {
			return false;
		}
	}
</script>

<section class="rounded-lg border border-border bg-surface p-6 shadow-sm sm:p-8">
	<div class="max-w-3xl">
		<p class="text-sm font-medium uppercase tracking-[0.18em] text-primary">Create source</p>
		<h1 class="mt-4 text-3xl font-semibold tracking-tight text-fg">Add a webhook destination</h1>
		<p class="mt-3 text-sm leading-6 text-muted sm:text-base">
			Create a unique ingest endpoint and define where captured webhook traffic should be forwarded.
		</p>
	</div>

	<form onsubmit={handleSubmit} class="mt-8 flex flex-col gap-6">
		<div class="flex flex-col gap-2">
			<label class="text-sm font-medium text-fg">
				Egress URL
				<input
					type="text"
					bind:value={data.egress_url}
					placeholder="https://example.com/egress"
					required
					disabled={loading}
					class="mt-2 w-full rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
				/>
			</label>
			<p class="text-sm text-muted">Use an `http` or `https` endpoint that should receive forwarded requests.</p>
			{#if egressError}
				<p class="text-sm text-error">{egressError}</p>
			{/if}
		</div>

		{#if data.static_headers}
			<div class="flex flex-col gap-3 rounded-lg border border-border-muted bg-elevated p-4">
				<div>
					<label class="text-sm font-medium text-fg">Static headers</label>
					<p class="mt-1 text-sm text-muted">
						Attach fixed headers to every forwarded request for this source.
					</p>
				</div>
				<InputMap bind:map={data.static_headers} disabled={loading} />
			</div>
		{/if}

		<div class="flex flex-col gap-2">
			<label class="text-sm font-medium text-fg">
				Description
				<textarea
					bind:value={data.description}
					placeholder="Optional description"
					disabled={loading}
					rows="4"
					class="mt-2 min-h-28 w-full rounded-md border border-border bg-elevated px-4 py-3 text-sm text-fg shadow-sm outline-none placeholder:text-subtle focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
				></textarea>
			</label>
		</div>

		<div class="flex flex-col gap-3 border-t border-border-muted pt-6 sm:flex-row sm:items-center sm:justify-between">
			<p class="text-sm text-muted">
				The source will be created immediately and appear in the sources list below.
			</p>
			<Button type="submit" disabled={loading || Boolean(egressError)}>
				{loading ? 'Creating Source...' : 'Create New Source'}
			</Button>
		</div>

		{#if error}
			<div class="rounded-md border border-error bg-surface px-4 py-3 text-sm text-error">
				Error: {error}
			</div>
		{/if}
	</form>
</section>
