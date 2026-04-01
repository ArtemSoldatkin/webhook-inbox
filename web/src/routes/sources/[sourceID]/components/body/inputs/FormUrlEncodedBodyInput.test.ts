import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import FormUrlEncodedBodyInputHost from '../../../../../../test/mocks/FormUrlEncodedBodyInputHost.svelte';

describe('FormUrlEncodedBodyInput', () => {
	it('serializes added key-value pairs into a url-encoded body', async () => {
		render(FormUrlEncodedBodyInputHost);

		await fireEvent.input(screen.getByLabelText('Key'), {
			target: { value: 'event' }
		});
		await fireEvent.input(screen.getByLabelText('Value'), {
			target: { value: 'payment_succeeded' }
		});
		await fireEvent.click(screen.getByRole('button', { name: 'Add key value pair' }));

		await waitFor(() => {
			expect(screen.getByTestId('body-state')).toHaveTextContent('event=payment_succeeded');
		});
		expect(screen.getByRole('button', { name: 'Clear' })).toBeEnabled();
	});

	it('url-encodes special characters in keys and values', async () => {
		render(FormUrlEncodedBodyInputHost);

		await fireEvent.input(screen.getByLabelText('Key'), {
			target: { value: 'user name' }
		});
		await fireEvent.input(screen.getByLabelText('Value'), {
			target: { value: 'a+b@example.com' }
		});
		await fireEvent.click(screen.getByRole('button', { name: 'Add key value pair' }));

		await waitFor(() => {
			expect(screen.getByTestId('body-state')).toHaveTextContent('user%20name=a%2Bb%40example.com');
		});
	});

	it('clears the encoded body when Clear is clicked', async () => {
		render(FormUrlEncodedBodyInputHost);

		await fireEvent.input(screen.getByLabelText('Key'), {
			target: { value: 'status' }
		});
		await fireEvent.input(screen.getByLabelText('Value'), {
			target: { value: 'ok' }
		});
		await fireEvent.click(screen.getByRole('button', { name: 'Add key value pair' }));

		await waitFor(() => {
			expect(screen.getByTestId('body-state')).toHaveTextContent('status=ok');
		});

		await fireEvent.click(screen.getByRole('button', { name: 'Clear' }));

		expect(screen.getByTestId('body-state')).toHaveTextContent('');
		expect(screen.getByTestId('error-state')).toHaveTextContent('');
		expect(screen.getByRole('button', { name: 'Clear' })).toBeDisabled();
	});
});
