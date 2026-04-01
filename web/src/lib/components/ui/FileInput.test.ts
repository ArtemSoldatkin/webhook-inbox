import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import FileInput from './FileInput.svelte';

describe('FileInput', () => {
	it('renders a file input with the provided attributes', () => {
		render(FileInput, {
			props: {
				'aria-label': 'Upload file',
				accept: '.json',
				disabled: true
			}
		});

		const input = screen.getByLabelText('Upload file');
		expect(input).toHaveAttribute('type', 'file');
		expect(input).toHaveAttribute('accept', '.json');
		expect(input).toBeDisabled();
	});
});
