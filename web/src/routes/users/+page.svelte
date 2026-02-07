<script lang="ts">
	import type { User } from '$lib/types';

	let email: string = '';
	let data: User | null = null;
	let loading = false;
	let error: string | null = null;

	async function createUser() {
		loading = true;
		error = null;
		try {
			const res = await fetch('/api/users', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					email
				})
			});
			if (!res.ok) {
				throw new Error('Failed to create user');
			}
			data = await res.json();
			email = '';
		} catch (err: unknown) {
			error = err instanceof Error ? err.message : String(err);
			console.error('Error creating user:', err);
		} finally {
			loading = false;
		}
	}
</script>

<h2>Create user</h2>
<form on:submit|preventDefault={createUser}>
	<label for="email">Email:</label>
	<input
		type="email"
		name="email"
		placeholder="Enter email"
		required
		bind:value={email}
		disabled={loading}
	/>
	<button type="submit" disabled={loading}>Create User</button>
</form>
{#if loading}
	<p>Loading...</p>
{:else if error}
	<p>Error: {error}</p>
{:else if data}
	<div>
		<h2>User ID: {data.ID}</h2>
		<p>Email: {data.Email}</p>
		<p>Created At: {new Date(data.CreatedAt).toLocaleString()}</p>
	</div>
{/if}
