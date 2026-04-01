import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';
import TestWebhook from './TestWebhook.svelte';

describe('TestWebhook', () => {
	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('submits a default GET request with static headers and no body', async () => {
		const fetchMock = vi.fn().mockResolvedValue({
			ok: true,
			statusText: 'OK'
		});
		vi.stubGlobal('fetch', fetchMock);

		render(TestWebhook, {
			props: {
				publicID: 'src_123',
				staticHeaders: { Authorization: 'Bearer token' }
			}
		});

		expect(screen.queryByText('Request body')).not.toBeInTheDocument();

		await fireEvent.submit(screen.getByRole('form', { name: 'Test webhook request' }));

		await waitFor(() => {
			expect(fetchMock).toHaveBeenCalledTimes(1);
		});

		const [url, options] = fetchMock.mock.calls[0];
		expect(url).toBe('/api/ingest/src_123');
		expect(options).toMatchObject({
			method: 'GET',
			headers: { Authorization: 'Bearer token' }
		});
		expect(options).not.toHaveProperty('body');
	});

	it('shows the body composer for POST and sends content type plus body', async () => {
		const fetchMock = vi.fn().mockResolvedValue({
			ok: true,
			statusText: 'OK'
		});
		vi.stubGlobal('fetch', fetchMock);

		render(TestWebhook, {
			props: {
				publicID: 'src_123',
				staticHeaders: { Authorization: 'Bearer token' }
			}
		});

		await fireEvent.change(screen.getByLabelText('HTTP Method'), {
			target: { value: 'POST' }
		});

		expect(screen.getByText('Request body')).toBeInTheDocument();
		expect(screen.getByLabelText('Content type')).toHaveValue('application/json');

		await fireEvent.submit(screen.getByRole('form', { name: 'Test webhook request' }));

		await waitFor(() => {
			expect(fetchMock).toHaveBeenCalledTimes(1);
		});

		const [url, options] = fetchMock.mock.calls[0];
		expect(url).toBe('/api/ingest/src_123');
		expect(options).toMatchObject({
			method: 'POST',
			headers: {
				Authorization: 'Bearer token',
				'Content-Type': 'application/json'
			},
			body: ''
		});
	});

	it('shows an error alert when the webhook request fails', async () => {
		const fetchMock = vi.fn().mockResolvedValue({
			ok: false,
			statusText: 'Bad Request'
		});
		const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});
		vi.stubGlobal('fetch', fetchMock);

		render(TestWebhook, {
			props: {
				publicID: 'src_123'
			}
		});

		await fireEvent.submit(screen.getByRole('form', { name: 'Test webhook request' }));

		expect(await screen.findByText('Failed to test webhook: Bad Request')).toBeInTheDocument();
		expect(consoleError).toHaveBeenCalled();
	});
});
