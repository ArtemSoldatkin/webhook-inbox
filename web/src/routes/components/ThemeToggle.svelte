<script lang="ts">
	import Select from '$lib/components/ui/Select.svelte';
	import { onMount } from 'svelte';

	/**
	 * Defines the possible theme options for the application.
	 * The 'light' and 'dark' options allow users to explicitly choose a theme,
	 * while the 'system' option lets the application follow the user's operating system preference.
	 */
	type Theme = 'light' | 'dark' | 'system';

	/** The 'mounted' state variable is used to determine if the component has been mounted on the client side. */
	let mounted = $state(false);
	/**
	 * The theme state is initialized to 'system', which means it will follow the user's system preference.
	 * When the user selects a different theme, it will be saved to localStorage and applied immediately.
	 */
	let theme = $state<Theme>('system');

	onMount(() => {
		const saved = localStorage.getItem('theme') as Theme | null;
		theme = saved ?? 'system';
		mounted = true;
	});

	$effect(() => {
		if (theme === 'system') {
			localStorage.removeItem('theme');
		} else {
			localStorage.setItem('theme', theme);
		}
		document.documentElement.classList.toggle(
			'dark',
			localStorage.theme === 'dark' ||
				(!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)
		);
	});
</script>

{#if mounted}
	<Select
		bind:value={theme}
		aria-label="Theme"
		options={[
			{ value: 'light', label: 'Light' },
			{ value: 'dark', label: 'Dark' },
			{ value: 'system', label: 'System' }
		]}
	/>
{/if}
