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

import JSONBodyView from './JSONBodyView.svelte';

describe('JSONBodyView', () => {
	it('renders formatted JSON without an error for valid input', () => {
		render(JSONBodyView, {
			props: {
				body: '{"user":"alice","active":true}'
			}
		});

		expect(screen.queryByText('Invalid JSON')).not.toBeInTheDocument();
		expect(screen.getByText(/"user"/)).toBeInTheDocument();
		expect(screen.getByText(/"alice"/)).toBeInTheDocument();
		expect(screen.getByText(/"active"/)).toBeInTheDocument();
	});

	it('shows an error and renders the raw input for invalid JSON', () => {
		render(JSONBodyView, {
			props: {
				body: '{"user":'
			}
		});

		expect(screen.getByText('Invalid JSON')).toBeInTheDocument();
		expect(screen.getByText('{"user":')).toBeInTheDocument();
	});
});
