import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import Link from './Link.svelte';

describe('Link', () => {
	it('renders an anchor with the requested href and variant classes', () => {
		render(Link, {
			props: {
				href: '/sources',
				variant: 'inline',
				'aria-label': 'Open sources'
			}
		});

		const link = screen.getByRole('link', { name: 'Open sources' });
		expect(link).toHaveAttribute('href', '/sources');
		expect(link.className).toContain('underline-offset-4');
	});
});
