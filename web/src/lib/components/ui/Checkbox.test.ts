import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import Checkbox from './Checkbox.svelte';

describe('Checkbox', () => {
	it('renders a checkbox and tracks checked state', async () => {
		render(Checkbox, {
			props: {
				'aria-label': 'Subscribe',
				value: false
			}
		});

		const checkbox = screen.getByRole('checkbox', { name: 'Subscribe' });
		expect(checkbox).not.toBeChecked();

		await fireEvent.click(checkbox);

		expect(checkbox).toBeChecked();
	});
});
