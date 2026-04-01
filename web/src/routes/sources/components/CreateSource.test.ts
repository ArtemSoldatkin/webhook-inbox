import { fireEvent, render, screen } from '@testing-library/svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';
import CreateSource from './CreateSource.svelte';

describe('CreateSource', () => {
	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('shows a validation error for an invalid egress url', async () => {
		render(CreateSource);

		const input = screen.getByLabelText('Egress URL');
		await fireEvent.input(input, {
			target: { value: 'invalid-url' }
		});

		expect(screen.getByText('Valid Egress URL is required')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Create New Source' })).toBeDisabled();
	});

	it('submits a valid source and resets the form', async () => {
		const fetchMock = vi.fn().mockResolvedValue({
			ok: true,
			json: async () => ({ id: 1 })
		});
		vi.stubGlobal('fetch', fetchMock);

		render(CreateSource);

		const urlInput = screen.getByLabelText('Egress URL');
		const descriptionInput = screen.getByLabelText('Description');

		await fireEvent.input(urlInput, {
			target: { value: 'https://example.com/webhook' }
		});
		await fireEvent.input(descriptionInput, {
			target: { value: 'Customer events' }
		});
		await fireEvent.submit(screen.getByRole('form', { name: 'Create source' }));

		expect(fetchMock).toHaveBeenCalledWith(
			'/api/sources',
			expect.objectContaining({ method: 'POST' })
		);
		expect(urlInput).toHaveValue('');
		expect(descriptionInput).toHaveValue('');
	});
});
