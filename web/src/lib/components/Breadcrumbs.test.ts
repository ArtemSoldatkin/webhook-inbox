import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import Breadcrumbs from './Breadcrumbs.svelte';

describe('Breadcrumbs', () => {
	it('renders links, active items, and separators', () => {
		render(Breadcrumbs, {
			props: {
				items: [
					{ label: 'Sources', href: '/sources' },
					{ label: 'Details', active: true }
				]
			}
		});

		expect(screen.getByRole('navigation', { name: 'Breadcrumb' })).toBeInTheDocument();
		expect(screen.getByRole('link', { name: 'Sources' })).toHaveAttribute('href', '/sources');
		expect(screen.getByText('Details')).toHaveAttribute('aria-current', 'page');
		expect(screen.getAllByText('/')).toHaveLength(1);
	});
});
