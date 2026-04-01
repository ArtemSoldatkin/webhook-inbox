import { fireEvent, render, screen } from '@testing-library/svelte';
import { createRawSnippet } from 'svelte';
import { describe, expect, it, vi } from 'vitest';
import Layout from './+layout.svelte';

vi.mock('$app/paths', () => ({
	resolve: (path: string) => path
}));

describe('routes/+layout', () => {
	it('renders the shared navigation and child content', () => {
		render(Layout, {
			props: {
				children: createRawSnippet(() => ({
					render: () => '<div>Route content</div>'
				}))
			}
		});

		expect(screen.getByRole('link', { name: 'Webhook Inbox' })).toHaveAttribute('href', '/');
		expect(screen.getByRole('navigation', { name: 'Primary' })).toBeInTheDocument();
		expect(screen.getByText('Route content')).toBeInTheDocument();
	});

	it('moves focus to the main content when skip link is activated', async () => {
		render(Layout, {
			props: {
				children: createRawSnippet(() => ({
					render: () => '<div>Route content</div>'
				}))
			}
		});

		const main = screen.getByRole('main');
		const scrollIntoView = vi.fn();
		main.scrollIntoView = scrollIntoView;

		await fireEvent.click(screen.getByRole('button', { name: 'Skip to main content' }));

		expect(document.activeElement).toBe(main);
		expect(scrollIntoView).toHaveBeenCalledWith({ block: 'start' });
	});
});
