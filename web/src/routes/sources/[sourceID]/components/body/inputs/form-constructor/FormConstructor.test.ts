import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import FormConstructor from './FormConstructor.svelte';

describe('FormConstructor', () => {
	it('shows an empty state when no fields exist', () => {
		render(FormConstructor, {
			props: {
				fields: []
			}
		});

		expect(
			screen.getByText('No fields added yet. Click "Add field" to start building your form.')
		).toBeInTheDocument();
	});

	it('adds a new empty text field', async () => {
		render(FormConstructor, {
			props: {
				fields: []
			}
		});

		await fireEvent.click(screen.getByRole('button', { name: 'Add field' }));

		expect(screen.getByRole('combobox')).toHaveValue('text');
		expect(screen.getByPlaceholderText('Enter name')).toHaveValue('');
		expect(screen.getByText('Name is required.')).toBeInTheDocument();
	});
});
