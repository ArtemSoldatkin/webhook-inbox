import { render, screen } from '@testing-library/svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';
import BodyView from './BodyView.svelte';

vi.mock('./body/views/ByteBodyView.svelte', async () => ({
	default: (await import('../../../../test/mocks/BodyViewChildStub.svelte')).default
}));

vi.mock('./body/views/FormUrlEncodedBodyView.svelte', async () => ({
	default: (await import('../../../../test/mocks/BodyViewChildStub.svelte')).default
}));

vi.mock('./body/views/JSONBodyView.svelte', async () => ({
	default: (await import('../../../../test/mocks/BodyViewChildStub.svelte')).default
}));

vi.mock('./body/views/PlainTextBodyView.svelte', async () => ({
	default: (await import('../../../../test/mocks/BodyViewChildStub.svelte')).default
}));

vi.mock('./body/views/XMLBodyView.svelte', async () => ({
	default: (await import('../../../../test/mocks/BodyViewChildStub.svelte')).default
}));

describe('BodyView', () => {
	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('shows a message when no body is provided', () => {
		render(BodyView);

		expect(screen.getByText('No body provided, cannot display content')).toBeInTheDocument();
	});

	it('shows a message when the body is empty', () => {
		render(BodyView, {
			props: {
				body: ''
			}
		});

		expect(screen.getByText('Body content is empty')).toBeInTheDocument();
	});

	it('renders decoded plain text through the plain text body view', () => {
		render(BodyView, {
			props: {
				body: btoa('hello world'),
				contentType: 'text/plain'
			}
		});

		expect(screen.getByTestId('body-view-child')).toHaveTextContent('hello world');
	});

	it('renders binary data through the byte body view', () => {
		render(BodyView, {
			props: {
				body: btoa('binary-data'),
				contentType: 'application/octet-stream'
			}
		});

		expect(screen.getByTestId('body-view-child')).toHaveTextContent('binary-data');
		expect(screen.getByTestId('body-view-child')).toHaveTextContent('application/octet-stream');
	});

	it('shows the unknown content type message when body parsing succeeds without a content type', () => {
		render(BodyView, {
			props: {
				body: btoa('payload')
			}
		});

		expect(screen.getByText('Content type unknown, cannot display body.')).toBeInTheDocument();
	});

	it('falls back to the raw body and shows a warning when base64 parsing fails', () => {
		const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});

		render(BodyView, {
			props: {
				body: '%%%not-base64%%%',
				contentType: 'application/json'
			}
		});

		expect(
			screen.getByText('Failed to parse body as Base64 string, falling back to raw content')
		).toBeInTheDocument();
		expect(screen.getByText('Original body: %%%not-base64%%%')).toBeInTheDocument();
		expect(consoleError).toHaveBeenCalled();
	});
});
