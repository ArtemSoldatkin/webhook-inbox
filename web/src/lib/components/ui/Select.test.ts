import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import Select from './Select.svelte';

describe('Select', () => {
	it('renders options and updates the selected value', async () => {
		render(Select, {
			props: {
				'aria-label': 'Status',
				value: 'active',
				options: [
					{ value: 'active', label: 'Active' },
					{ value: 'paused', label: 'Paused' }
				]
			}
		});

		const select = screen.getByLabelText('Status');
		expect(screen.getAllByRole('option')).toHaveLength(2);
		expect(select).toHaveValue('active');

		await fireEvent.change(select, {
			target: { value: 'paused' }
		});

		expect(select).toHaveValue('paused');
	});
});
