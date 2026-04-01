import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import KeyValueList from './KeyValueList.svelte';

describe('KeyValueList', () => {
	it('renders item labels and maps null values to N/A', () => {
		render(KeyValueList, {
			props: {
				items: [
					{ label: 'Content-Type', value: 'application/json' },
					{ label: 'Attempt', value: null }
				]
			}
		});

		expect(screen.getByText('Content-Type')).toBeInTheDocument();
		expect(screen.getByText('application/json')).toBeInTheDocument();
		expect(screen.getByText('Attempt')).toBeInTheDocument();
		expect(screen.getByText('N/A')).toBeInTheDocument();
	});

	it('renders the empty state text when no items are present', () => {
		render(KeyValueList, {
			props: {
				items: [],
				emptyStateText: 'Nothing to show'
			}
		});

		expect(screen.getByText('Nothing to show')).toBeInTheDocument();
	});
});
