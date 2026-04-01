import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import SectionHeader from './SectionHeader.svelte';

describe('SectionHeader', () => {
	it('renders eyebrow, title, and description', () => {
		render(SectionHeader, {
			props: {
				eyebrow: 'Overview',
				title: 'Delivery Attempts',
				description: 'Inspect recent delivery state changes.',
				titleAs: 'h3'
			}
		});

		expect(screen.getByText('Overview')).toBeInTheDocument();
		expect(
			screen.getByRole('heading', { level: 3, name: 'Delivery Attempts' })
		).toBeInTheDocument();
		expect(screen.getByText('Inspect recent delivery state changes.')).toBeInTheDocument();
	});
});
