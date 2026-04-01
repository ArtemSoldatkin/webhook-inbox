import { render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import SourcesPage from './+page.svelte';

vi.mock('$app/paths', () => ({
	resolve: (path: string) => path
}));

vi.mock('$lib/api', () => ({
	fetchPaginatedData: vi.fn().mockResolvedValue({
		data: [],
		next_cursor: null,
		has_next: false
	})
}));

describe('routes/sources/+page', () => {
	it('renders the create source form and sources list state', async () => {
		render(SourcesPage);

		expect(screen.getByRole('form', { name: 'Create source' })).toBeInTheDocument();
		expect(
			screen.getByRole('heading', { name: 'Manage registered endpoints' })
		).toBeInTheDocument();
		expect(await screen.findByText('No sources found.')).toBeInTheDocument();
	});
});
