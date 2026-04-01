import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import FormUrlEncodedBodyView from './FormUrlEncodedBodyView.svelte';

describe('FormUrlEncodedBodyView', () => {
	it('renders parsed key value pairs from the encoded body', () => {
		render(FormUrlEncodedBodyView, {
			props: {
				body: 'name=Alice+Smith&role=admin%2Feditor'
			}
		});

		expect(screen.getByText('name')).toBeInTheDocument();
		expect(screen.getByText('Alice Smith')).toBeInTheDocument();
		expect(screen.getByText('role')).toBeInTheDocument();
		expect(screen.getByText('admin/editor')).toBeInTheDocument();
	});

	it('renders the empty state when no form values are present', () => {
		render(FormUrlEncodedBodyView, {
			props: {
				body: ''
			}
		});

		expect(screen.getByText('No form values provided.')).toBeInTheDocument();
	});
});
