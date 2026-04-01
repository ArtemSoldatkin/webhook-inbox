import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import InputConstructorHost from '../../../../../../../test/mocks/InputConstructorHost.svelte';

describe('InputConstructor', () => {
	it('shows a validation warning when the field name is empty', () => {
		render(InputConstructorHost, {
			props: {
				initialField: {
					type: 'text',
					name: '',
					value: 'hello'
				}
			}
		});

		expect(screen.getByText('Name is required.')).toBeInTheDocument();
	});

	it('switches input type and resets the value for number fields', async () => {
		render(InputConstructorHost, {
			props: {
				initialField: {
					type: 'text',
					name: 'age',
					value: '123'
				}
			}
		});

		await fireEvent.change(screen.getByRole('combobox'), {
			target: { value: 'number' }
		});

		const numberInput = screen.getByPlaceholderText('Enter number');
		expect(numberInput).toHaveValue(null);
		expect(screen.getByTestId('field-state')).toHaveTextContent('"type":"number"');
		expect(screen.getByTestId('field-state')).toHaveTextContent('"value":null');
	});

	it('renders a checkbox editor and clears its value', async () => {
		render(InputConstructorHost, {
			props: {
				initialField: {
					type: 'checkbox',
					name: 'enabled',
					value: true
				}
			}
		});

		const checkbox = screen.getByRole('checkbox');
		expect(checkbox).toBeChecked();

		await fireEvent.click(screen.getByRole('button', { name: 'Clear field value' }));

		expect(checkbox).not.toBeChecked();
		expect(screen.getByTestId('field-state')).toHaveTextContent('"value":false');
	});

	it('clears text field values', async () => {
		render(InputConstructorHost, {
			props: {
				initialField: {
					type: 'text',
					name: 'title',
					value: 'hello'
				}
			}
		});

		const textInput = screen.getByPlaceholderText('Enter text');
		expect(textInput).toHaveValue('hello');

		await fireEvent.click(screen.getByRole('button', { name: 'Clear field value' }));

		expect(textInput).toHaveValue('');
		expect(screen.getByTestId('field-state')).toHaveTextContent('"value":""');
	});
});
