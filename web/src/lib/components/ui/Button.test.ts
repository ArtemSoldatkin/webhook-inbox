import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import Button from './Button.svelte';

describe('Button', () => {
	it('defaults to type button and applies the selected variant class', () => {
		render(Button, {
			props: {
				variant: 'secondary'
			}
		});

		const button = screen.getByRole('button');
		expect(button).toHaveAttribute('type', 'button');
		expect(button.className).toContain('bg-surface');
	});
});
