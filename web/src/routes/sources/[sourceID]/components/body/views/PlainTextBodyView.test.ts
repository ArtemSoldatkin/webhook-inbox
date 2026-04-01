import { render } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import PlainTextBodyView from './PlainTextBodyView.svelte';

describe('PlainTextBodyView', () => {
	it('renders the provided plain text body verbatim', () => {
		const { container } = render(PlainTextBodyView, {
			props: {
				body: 'line 1\nline 2'
			}
		});

		const pre = container.querySelector('pre');

		expect(pre).not.toBeNull();
		expect(pre?.textContent).toBe('line 1\nline 2');
	});
});
