import { render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import HomePage from './+page.svelte';

vi.mock('$app/paths', () => ({
	resolve: (path: string) => path
}));

describe('routes/+page', () => {
	it('renders the landing page hero and quick start content', () => {
		const { container } = render(HomePage);

		expect(
			screen.getByRole('heading', {
				level: 1,
				name: 'See every webhook the way your provider actually sent it.'
			})
		).toBeInTheDocument();
		expect(screen.getByRole('link', { name: 'Open Sources' })).toHaveAttribute('href', '/sources');
		expect(screen.getByText('Typical use cases')).toBeInTheDocument();
		expect(screen.getByText('01')).toBeInTheDocument();
		expect(screen.getByText('02')).toBeInTheDocument();
		expect(screen.getByText('03')).toBeInTheDocument();
		expect(container.textContent).toContain('curl -X POST');
		expect(container.textContent).toContain('/api/v1/ingest/YOUR_SOURCE_KEY');
	});
});
