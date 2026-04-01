import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import XMLBodyInputHost from '../../../../../../test/mocks/XMLBodyInputHost.svelte';

describe('XMLBodyInput', () => {
	it('flags invalid XML and disables formatting', async () => {
		render(XMLBodyInputHost, {
			props: {
				initialBody: '<root><child></root>',
				initialError: null
			}
		});

		await waitFor(() => {
			expect(screen.getByTestId('error-state')).toHaveTextContent('Invalid XML format');
		});
		expect(screen.getByRole('button', { name: 'Format' })).toBeDisabled();
	});

	it('formats valid XML into a pretty-printed body', async () => {
		render(XMLBodyInputHost, {
			props: {
				initialBody: '<root><child>value</child></root>',
				initialError: null
			}
		});

		await fireEvent.click(screen.getByRole('button', { name: 'Format' }));

		const formattedBody = screen.getByTestId('body-state').textContent?.replaceAll('\r\n', '\n');
		const textarea = screen.getByLabelText('XML request body') as HTMLTextAreaElement;

		expect(formattedBody).toBe('<root>\n    <child>\n        value\n    </child>\n</root>');
		expect(textarea.value.replaceAll('\r\n', '\n')).toBe(
			'<root>\n    <child>\n        value\n    </child>\n</root>'
		);
		expect(screen.getByTestId('error-state').textContent).toBe('');
	});

	it('clears the body and resets the validation state', async () => {
		render(XMLBodyInputHost, {
			props: {
				initialBody: '<root><child>value</child></root>',
				initialError: null
			}
		});

		await fireEvent.click(screen.getByRole('button', { name: 'Clear' }));

		expect(screen.getByLabelText('XML request body')).toHaveValue('');
		expect(screen.getByTestId('body-state').textContent).toBe('');
		expect(screen.getByTestId('error-state').textContent).toBe('');
		expect(screen.getByRole('button', { name: 'Clear' })).toBeDisabled();
	});
});
