import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import BodyInput from './BodyInput.svelte';

vi.mock('./body/inputs/ByteBodyInput.svelte', async () => ({
	default: (await import('../../../../test/mocks/BodyInputChildStub.svelte')).default
}));

vi.mock('./body/inputs/FormDataBodyInput.svelte', async () => ({
	default: (await import('../../../../test/mocks/BodyInputChildStub.svelte')).default
}));

vi.mock('./body/inputs/FormUrlEncodedBodyInput.svelte', async () => ({
	default: (await import('../../../../test/mocks/BodyInputChildStub.svelte')).default
}));

vi.mock('./body/inputs/JSONBodyInput.svelte', async () => ({
	default: (await import('../../../../test/mocks/BodyInputChildStub.svelte')).default
}));

vi.mock('./body/inputs/PlainTextBodyInput.svelte', async () => ({
	default: (await import('../../../../test/mocks/BodyInputChildStub.svelte')).default
}));

vi.mock('./body/inputs/XMLBodyInput.svelte', async () => ({
	default: (await import('../../../../test/mocks/BodyInputChildStub.svelte')).default
}));

describe('BodyInput', () => {
	it('renders the JSON editor by default and clears existing text body on mount', () => {
		render(BodyInput, {
			props: {
				textBody: '{"ok":true}',
				formDataBody: new FormData(),
				contentType: 'application/json'
			}
		});

		expect(screen.getByLabelText('Content type')).toHaveValue('application/json');
		expect(screen.getByTestId('body-input-child')).toBeInTheDocument();
		expect(screen.getByTestId('body-input-child')).not.toHaveTextContent('{"ok":true}');
	});

	it('resets body state when the content type changes and renders the new editor', async () => {
		render(BodyInput, {
			props: {
				textBody: '{"ok":true}',
				formDataBody: new FormData(),
				contentType: 'application/json'
			}
		});

		await fireEvent.change(screen.getByLabelText('Content type'), {
			target: { value: 'multipart/form-data' }
		});

		expect(screen.getByLabelText('Content type')).toHaveValue('multipart/form-data');
		expect(screen.getByTestId('body-input-child')).toHaveTextContent('form-data');
		expect(screen.getByTestId('body-input-child')).not.toHaveTextContent('{"ok":true}');
	});

	it('shows the binary editor when application/octet-stream is selected', async () => {
		render(BodyInput, {
			props: {
				textBody: '',
				formDataBody: new FormData(),
				contentType: 'text/plain'
			}
		});

		await fireEvent.change(screen.getByLabelText('Content type'), {
			target: { value: 'application/octet-stream' }
		});

		expect(screen.getByLabelText('Content type')).toHaveValue('application/octet-stream');
		expect(screen.getByTestId('body-input-child')).toBeInTheDocument();
	});
});
