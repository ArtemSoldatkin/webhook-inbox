import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import FormDataBodyInputHost from '../../../../../../test/mocks/FormDataBodyInputHost.svelte';

vi.mock('./form-constructor/FormConstructor.svelte', async () => ({
	default: (await import('../../../../../../test/mocks/FormConstructorStub.svelte')).default
}));

describe('FormDataBodyInput', () => {
	it('builds text and number fields into FormData', async () => {
		render(FormDataBodyInputHost);

		await fireEvent.click(screen.getByRole('button', { name: 'load-text-field' }));

		await waitFor(() => {
			expect(screen.getByTestId('formdata-state')).toHaveTextContent(
				'[["title","hello"],["count","3"]]'
			);
		});
		expect(screen.getByTestId('error-state')).toHaveTextContent('');
	});

	it('serializes checkbox fields as on/off values', async () => {
		render(FormDataBodyInputHost);

		await fireEvent.click(screen.getByRole('button', { name: 'load-checkbox-field' }));

		await waitFor(() => {
			expect(screen.getByTestId('formdata-state')).toHaveTextContent('[["enabled","on"]]');
		});
	});

	it('appends uploaded files to FormData', async () => {
		render(FormDataBodyInputHost);

		await fireEvent.click(screen.getByRole('button', { name: 'load-file-field' }));

		await waitFor(() => {
			expect(screen.getByTestId('formdata-state')).toHaveTextContent(
				'[["attachment","file:demo.txt"]]'
			);
		});
	});
});
