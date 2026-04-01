import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import FilterBar from './FilterBar.svelte';

describe('FilterBar', () => {
	it('renders search, filter, and toggles sort direction state', async () => {
		render(FilterBar, {
			props: {
				searchQuery: '',
				sortDirection: 'ASC',
				filterName: 'status',
				filterOptions: ['active', 'paused'],
				filter: '*'
			}
		});

		expect(screen.getByLabelText('Search')).toBeInTheDocument();
		expect(screen.getByLabelText('Filter by status')).toBeInTheDocument();
		expect(screen.getByRole('option', { name: 'All' })).toBeInTheDocument();

		const sortButton = screen.getByRole('button', {
			name: 'Sort ascending. Activate to sort descending'
		});
		await fireEvent.click(sortButton);

		expect(
			screen.getByRole('button', {
				name: 'Sort descending. Activate to sort ascending'
			})
		).toHaveAttribute('aria-pressed', 'true');
	});
});
