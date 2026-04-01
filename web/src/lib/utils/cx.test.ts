import { describe, expect, it } from 'vitest';
import { cx } from './cx';

describe('cx', () => {
	it('filters out falsy class values', () => {
		expect(cx('px-4', false, null, undefined, '', 'py-2')).toBe('px-4 py-2');
	});

	it('lets later Tailwind utility classes win', () => {
		expect(cx('px-2 text-sm', 'px-4', 'text-base')).toBe('px-4 text-base');
	});

	it('preserves non-conflicting classes while merging conflicts', () => {
		expect(cx('rounded-md px-2', 'bg-surface', 'px-6')).toBe('rounded-md bg-surface px-6');
	});
});
