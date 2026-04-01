import { render } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import Eyebrow from './Eyebrow.svelte';

describe('Eyebrow', () => {
	it('renders as a label when requested', () => {
		const { container } = render(Eyebrow, {
			props: {
				as: 'label',
				variant: 'strong'
			}
		});

		const label = container.querySelector('label');
		expect(label).toBeInTheDocument();
		expect(label?.className).toContain('text-primary');
	});
});
