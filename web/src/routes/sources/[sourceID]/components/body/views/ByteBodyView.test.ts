import { fireEvent, render, screen } from '@testing-library/svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';
import ByteBodyView from './ByteBodyView.svelte';

describe('ByteBodyView', () => {
	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('renders the body size, hex preview, and disables download when empty', () => {
		render(ByteBodyView, {
			props: {
				body: '',
				contentType: 'application/octet-stream'
			}
		});

		expect(screen.getByText('0 bytes')).toBeInTheDocument();
		expect(screen.getByText('...')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Download as file' })).toBeDisabled();
	});

	it('downloads the body as a file when requested', async () => {
		const createObjectURL = vi.fn(() => 'blob:test-url');
		const revokeObjectURL = vi.fn();
		const clickSpy = vi.fn();
		const appendChildSpy = vi.spyOn(document.body, 'appendChild');
		const removeChildSpy = vi.spyOn(document.body, 'removeChild');
		const createElementSpy = vi.spyOn(document, 'createElement');

		vi.stubGlobal('URL', {
			createObjectURL,
			revokeObjectURL
		});
		createElementSpy.mockImplementation(((tagName: string) => {
			if (tagName === 'a') {
				const anchor = document.createElementNS('http://www.w3.org/1999/xhtml', 'a');
				anchor.click = clickSpy;
				return anchor;
			}

			return document.createElementNS('http://www.w3.org/1999/xhtml', tagName);
		}) as typeof document.createElement);

		render(ByteBodyView, {
			props: {
				body: 'ABC',
				contentType: 'application/octet-stream'
			}
		});

		expect(screen.getByText('3 bytes')).toBeInTheDocument();
		expect(screen.getByText('41 42 43...')).toBeInTheDocument();

		await fireEvent.click(screen.getByRole('button', { name: 'Download as file' }));

		expect(createObjectURL).toHaveBeenCalledTimes(1);
		expect(clickSpy).toHaveBeenCalledTimes(1);
		expect(
			appendChildSpy.mock.calls.some(([node]) => (node as Element).tagName === 'A')
		).toBe(true);
		expect(
			removeChildSpy.mock.calls.some(([node]) => (node as Element).tagName === 'A')
		).toBe(true);
		expect(revokeObjectURL).toHaveBeenCalledWith('blob:test-url');
	});
});
