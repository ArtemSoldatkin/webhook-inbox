import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';
import ByteBodyInputHost from '../../../../../../test/mocks/ByteBodyInputHost.svelte';
import ByteBodyInput from './ByteBodyInput.svelte';

class FileReaderSuccessMock {
	result: ArrayBuffer | null = null;
	onload: null | (() => void) = null;
	onerror: null | (() => void) = null;

	readAsArrayBuffer(): void {
		this.result = Uint8Array.from([72, 101, 108, 108, 111]).buffer;
		this.onload?.();
	}
}

class FileReaderErrorMock {
	result: ArrayBuffer | null = null;
	onload: null | (() => void) = null;
	onerror: null | (() => void) = null;

	readAsArrayBuffer(): void {
		this.onerror?.();
	}
}

describe('ByteBodyInput', () => {
	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('accepts a valid base64 body and enables clearing', async () => {
		render(ByteBodyInput, {
			props: {
				body: '',
				error: null
			}
		});

		const textarea = screen.getByLabelText('Base64 request body');
		const clearButton = screen.getByRole('button', { name: 'Clear' });

		expect(clearButton).toBeDisabled();

		await fireEvent.input(textarea, {
			target: { value: 'SGVsbG8=' }
		});

		expect(textarea).toHaveValue('SGVsbG8=');
		expect(clearButton).toBeEnabled();
	});

	it('clears the body when the clear button is clicked', async () => {
		render(ByteBodyInput, {
			props: {
				body: 'SGVsbG8=',
				error: null
			}
		});

		const textarea = screen.getByLabelText('Base64 request body');
		await fireEvent.click(screen.getByRole('button', { name: 'Clear' }));

		expect(textarea).toHaveValue('');
		expect(screen.getByRole('button', { name: 'Clear' })).toBeDisabled();
	});

	it('flags invalid base64 input', async () => {
		render(ByteBodyInputHost, {
			props: {
				initialBody: '',
				initialError: null
			}
		});

		await fireEvent.input(screen.getByLabelText('Base64 request body'), {
			target: { value: '%%%not-base64%%%' }
		});

		await waitFor(() => {
			expect(screen.getByTestId('error-state')).toHaveTextContent('Invalid base64 string');
		});
	});

	it('reads an uploaded file and stores its content as base64', async () => {
		vi.stubGlobal('FileReader', FileReaderSuccessMock);

		render(ByteBodyInputHost, {
			props: {
				initialBody: '',
				initialError: null
			}
		});

		const input = screen.getByLabelText('Upload file');
		const file = new File(['Hello'], 'hello.txt', { type: 'text/plain' });

		await fireEvent.change(input, {
			target: { files: [file] }
		});

		expect(screen.getByLabelText('Base64 request body')).toHaveValue('SGVsbG8=');
		expect(screen.getByTestId('body-state')).toHaveTextContent('SGVsbG8=');
	});

	it('shows an error when reading an uploaded file fails', async () => {
		vi.stubGlobal('FileReader', FileReaderErrorMock);

		render(ByteBodyInputHost, {
			props: {
				initialBody: '',
				initialError: null
			}
		});

		const input = screen.getByLabelText('Upload file');
		const file = new File(['Hello'], 'hello.txt', { type: 'text/plain' });

		await fireEvent.change(input, {
			target: { files: [file] }
		});

		await waitFor(() => {
			expect(screen.getByTestId('error-state')).toHaveTextContent('Failed to read file');
		});
	});
});
