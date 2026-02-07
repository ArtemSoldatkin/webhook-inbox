<script lang="ts">
    import { page } from '$app/state';
    import type { Endpoint } from '$lib/types';
    import { onMount } from 'svelte';

    const {endpointID} = page.params;

    let data: Endpoint | null = null;
    let loading = false;
    let error: string | null = null;

    async function fetchEndpoint() {
        loading = true;
        error = null;
        try {
            const res = await fetch(`/api/endpoints/${endpointID}`);
            if (!res.ok) {
                throw new Error('Failed to fetch data');
            }
            data = await res.json();
        } catch (err: unknown) {
            error = err instanceof Error ? err.message : String(err);
            console.error('Error fetching endpoint:', err);
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        if (endpointID) {
            fetchEndpoint();
        } else {
            error = 'No endpoint ID provided in URL';
        }
    });

</script>

{#if loading}
	<p>Loading...</p>
{:else if error}
	<p>Error: {error}</p>
{:else if data}
    <div>
        <h2>{data.Name}</h2>
        <p>{data.Description}</p>
        <p>{data.Url}</p>
        {#each Object.entries(data.Headers) as [headerKey, headerValue]}
            <p>{headerKey}: {headerValue}</p>
        {/each}
        <p>Active: {data.IsActive ? 'Yes' : 'No'}</p>
        <p>Created At: {new Date(data.CreatedAt).toLocaleString()}</p>
    </div>
{/if}
<nav>
    <ul>
        <li><a href="/endpoints/{endpointID}/webhooks">Webhooks</a></li>
    </ul>
</nav>