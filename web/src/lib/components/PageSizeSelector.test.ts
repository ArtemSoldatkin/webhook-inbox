import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import PageSizeSelector from './PageSizeSelector.svelte';

describe('PageSizeSelector', () => {
	it('renders available page size options and updates the selected value', async () => {
		render(PageSizeSelector, {
			props: {
				pageSize: 20
			}
		});

		const select = screen.getByLabelText('Page size');
		expect(screen.getAllByRole('option')).toHaveLength(6);
		expect(select).toHaveValue('20');

		await fireEvent.change(select, {
			target: { value: '50' }
		});

		expect(select).toHaveValue('50');
	});
});
