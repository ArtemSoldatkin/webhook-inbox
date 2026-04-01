import { render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';

function escapeHtml(value: string): string {
	return value.replaceAll('&', '&amp;').replaceAll('<', '&lt;').replaceAll('>', '&gt;');
}

vi.mock('highlight.js/lib/core', () => ({
	default: {
		highlight: (code: string) => ({ value: escapeHtml(code) })
	}
}));

import XMLBodyView from './XMLBodyView.svelte';

describe('XMLBodyView', () => {
	it('renders formatted XML without an error for valid input', () => {
		const { container } = render(XMLBodyView, {
			props: {
				body: '<root><child>value</child></root>'
			}
		});

		const code = container.querySelector('code');

		expect(screen.queryByText('Invalid XML')).not.toBeInTheDocument();
		expect(code?.textContent?.replaceAll('\r\n', '\n')).toBe(
			'<root>\n    <child>\n        value\n    </child>\n</root>'
		);
	});

	it('shows an error and renders the raw input for invalid XML', () => {
		const { container } = render(XMLBodyView, {
			props: {
				body: '<root><'
			}
		});

		const code = container.querySelector('code');

		expect(screen.getByText('Invalid XML')).toBeInTheDocument();
		expect(code?.textContent).toBe('<root><');
	});
});
