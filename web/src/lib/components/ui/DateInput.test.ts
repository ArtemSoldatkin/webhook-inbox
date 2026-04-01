import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import DateInput from './DateInput.svelte';

describe('DateInput', () => {
	it('renders a date input and updates its value', async () => {
		render(DateInput, {
			props: {
				'aria-label': 'Delivery date',
				value: null
			}
		});

		const input = screen.getByLabelText('Delivery date');
		expect(input).toHaveAttribute('type', 'date');

		await fireEvent.input(input, {
			target: { value: '2026-04-01' }
		});

		expect(input).toHaveValue('2026-04-01');
	});
});
