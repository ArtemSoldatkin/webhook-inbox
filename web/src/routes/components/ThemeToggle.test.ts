import { fireEvent, render, screen } from '@testing-library/svelte';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import ThemeToggle from './ThemeToggle.svelte';

describe('ThemeToggle', () => {
	beforeEach(() => {
		localStorage.clear();
		document.documentElement.classList.remove('dark');
		Object.defineProperty(window, 'matchMedia', {
			writable: true,
			value: vi.fn().mockImplementation(() => ({
				matches: false,
				addEventListener: vi.fn(),
				removeEventListener: vi.fn()
			}))
		});
	});

	it('uses a saved theme from localStorage on mount', async () => {
		localStorage.setItem('theme', 'dark');

		render(ThemeToggle);

		const select = await screen.findByLabelText('Theme');
		expect(select).toHaveValue('dark');
		expect(document.documentElement.classList.contains('dark')).toBe(true);
	});

	it('updates the stored theme and clears dark mode when set to light', async () => {
		localStorage.setItem('theme', 'dark');

		render(ThemeToggle);

		const select = await screen.findByLabelText('Theme');
		await fireEvent.change(select, { target: { value: 'light' } });

		expect(localStorage.getItem('theme')).toBe('light');
		expect(document.documentElement.classList.contains('dark')).toBe(false);
	});
});
