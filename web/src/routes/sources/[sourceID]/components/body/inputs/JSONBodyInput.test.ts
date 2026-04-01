import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import JSONBodyInputHost from '../../../../../../test/mocks/JSONBodyInputHost.svelte';

describe('JSONBodyInput', () => {
	it('flags invalid JSON and disables formatting', async () => {
		render(JSONBodyInputHost, {
			props: {
				initialBody: '{"broken":',
				initialError: null
			}
		});

		await waitFor(() => {
			expect(screen.getByTestId('error-state')).toHaveTextContent('Invalid JSON format');
		});
		expect(screen.getByRole('button', { name: 'Format' })).toBeDisabled();
	});

	it('formats valid JSON into a pretty-printed body', async () => {
		render(JSONBodyInputHost, {
			props: {
				initialBody: '{"ok":true,"count":2}',
				initialError: null
			}
		});

		await fireEvent.click(screen.getByRole('button', { name: 'Format' }));

		expect(screen.getByTestId('body-state').textContent).toBe('{\n  "ok": true,\n  "count": 2\n}');
		expect(screen.getByLabelText('JSON request body')).toHaveValue(
			'{\n  "ok": true,\n  "count": 2\n}'
		);
		expect(screen.getByTestId('error-state')).toHaveTextContent('');
	});

	it('clears the body and revalidates the empty state', async () => {
		render(JSONBodyInputHost, {
			props: {
				initialBody: '{"ok":true}',
				initialError: null
			}
		});

		await fireEvent.click(screen.getByRole('button', { name: 'Clear' }));

		expect(screen.getByTestId('body-state').textContent).toBe('');
		expect(screen.getByTestId('error-state')).toHaveTextContent('Invalid JSON format');
		expect(screen.getByRole('button', { name: 'Clear' })).toBeDisabled();
	});
});
