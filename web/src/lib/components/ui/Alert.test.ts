import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import Alert from './Alert.svelte';

describe('Alert', () => {
	it('uses alert semantics for the error variant', () => {
		render(Alert, {
			props: {
				variant: 'error',
				title: 'Request failed'
			}
		});

		const alert = screen.getByRole('alert');
		expect(alert).toHaveAttribute('aria-live', 'assertive');
		expect(screen.getByText('Request failed')).toBeInTheDocument();
	});
});
