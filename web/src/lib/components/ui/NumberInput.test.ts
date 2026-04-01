import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import NumberInput from './NumberInput.svelte';

describe('NumberInput', () => {
	it('parses numeric changes and calls the provided onchange handler', async () => {
		const onchange = vi.fn();

		render(NumberInput, {
			props: {
				'aria-label': 'Retries',
				value: 1,
				onchange
			}
		});

		const input = screen.getByLabelText('Retries');
		await fireEvent.change(input, {
			target: { value: '42' }
		});

		expect(input).toHaveValue(42);
		expect(onchange).toHaveBeenCalledTimes(1);
	});

	it('clears to an empty string when the input is emptied', async () => {
		render(NumberInput, {
			props: {
				'aria-label': 'Retries',
				value: 5
			}
		});

		const input = screen.getByLabelText('Retries');
		await fireEvent.change(input, {
			target: { value: '' }
		});

		expect(input).toHaveValue(null);
	});
});
