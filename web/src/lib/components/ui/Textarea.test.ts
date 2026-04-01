import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import Textarea from './Textarea.svelte';

describe('Textarea', () => {
	it('renders with the provided row count and updates its value', async () => {
		render(Textarea, {
			props: {
				'aria-label': 'Payload',
				value: 'hello',
				rows: 4
			}
		});

		const textarea = screen.getByLabelText('Payload');
		expect(textarea).toHaveAttribute('rows', '4');

		await fireEvent.input(textarea, {
			target: { value: 'updated payload' }
		});

		expect(textarea).toHaveValue('updated payload');
	});
});
