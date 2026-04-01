import { render } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import Badge from './Badge.svelte';

describe('Badge', () => {
	it('renders as a div when requested and applies variant classes', () => {
		const { container } = render(Badge, {
			props: {
				as: 'div',
				variant: 'success',
				appearance: 'outline'
			}
		});

		const badge = container.querySelector('div');
		expect(badge).toBeInTheDocument();
		expect(badge?.className).toContain('border-success');
		expect(badge?.className).toContain('bg-transparent');
	});
});
