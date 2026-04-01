import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import TextInput from './TextInput.svelte';

describe('TextInput', () => {
	it('renders a text input and updates its value', async () => {
		render(TextInput, {
			props: {
				'aria-label': 'Name',
				value: 'Initial'
			}
		});

		const input = screen.getByLabelText('Name');
		expect(input).toHaveAttribute('type', 'text');

		await fireEvent.input(input, {
			target: { value: 'Updated' }
		});

		expect(input).toHaveValue('Updated');
	});
});
