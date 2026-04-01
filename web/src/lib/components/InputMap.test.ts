import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import InputMap from './InputMap.svelte';

describe('InputMap', () => {
	it('adds a new key value pair', async () => {
		render(InputMap, {
			props: {
				map: {}
			}
		});

		expect(screen.getByText('No static headers added yet.')).toBeInTheDocument();

		await fireEvent.input(screen.getByLabelText('Key'), {
			target: { value: 'Authorization' }
		});
		await fireEvent.input(screen.getByLabelText('Value'), {
			target: { value: 'Bearer token' }
		});
		await fireEvent.click(screen.getByRole('button', { name: 'Add key value pair' }));

		expect(screen.getByText('Authorization')).toBeInTheDocument();
		expect(screen.getByDisplayValue('Bearer token')).toBeInTheDocument();
		expect(screen.queryByText('No static headers added yet.')).not.toBeInTheDocument();
	});

	it('shows a validation error for duplicate keys', async () => {
		render(InputMap, {
			props: {
				map: {
					Authorization: 'Bearer token'
				}
			}
		});

		const keyInput = screen.getByLabelText('Key');
		await fireEvent.input(keyInput, {
			target: { value: 'Authorization' }
		});
		await fireEvent.blur(keyInput);

		expect(screen.getByText('Key already exists in the map.')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Add key value pair' })).toBeDisabled();
	});

	it('removes an existing key value pair', async () => {
		render(InputMap, {
			props: {
				map: {
					Authorization: 'Bearer token'
				}
			}
		});

		await fireEvent.click(screen.getByRole('button', { name: 'Remove Authorization' }));

		expect(screen.queryByText('Authorization')).not.toBeInTheDocument();
		expect(screen.getByText('No static headers added yet.')).toBeInTheDocument();
	});
});
