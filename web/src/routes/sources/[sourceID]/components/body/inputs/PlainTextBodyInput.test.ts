import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import PlainTextBodyInputHost from '../../../../../../test/mocks/PlainTextBodyInputHost.svelte';

describe('PlainTextBodyInput', () => {
	it('renders the bound plain text body and enables clearing', () => {
		render(PlainTextBodyInputHost, {
			props: {
				initialBody: 'hello world',
				initialError: null
			}
		});

		expect(screen.getByLabelText('Plain text request body')).toHaveValue('hello world');
		expect(screen.getByRole('button', { name: 'Clear' })).toBeEnabled();
		expect(screen.getByTestId('body-state')).toHaveTextContent('hello world');
	});

	it('updates the bound body when the user types', async () => {
		render(PlainTextBodyInputHost, {
			props: {
				initialBody: '',
				initialError: null
			}
		});

		await fireEvent.input(screen.getByLabelText('Plain text request body'), {
			target: { value: 'updated body' }
		});

		expect(screen.getByTestId('body-state')).toHaveTextContent('updated body');
		expect(screen.getByRole('button', { name: 'Clear' })).toBeEnabled();
	});

	it('clears the body and resets any existing error', async () => {
		render(PlainTextBodyInputHost, {
			props: {
				initialBody: 'keep me',
				initialError: 'Previous error'
			}
		});

		await fireEvent.click(screen.getByRole('button', { name: 'Clear' }));

		expect(screen.getByLabelText('Plain text request body')).toHaveValue('');
		expect(screen.getByTestId('body-state').textContent).toBe('');
		expect(screen.getByTestId('error-state').textContent).toBe('');
		expect(screen.getByRole('button', { name: 'Clear' })).toBeDisabled();
	});
});
