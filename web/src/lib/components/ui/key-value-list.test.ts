import { describe, expect, it } from 'vitest';
import { stringArrayRecordToKeyValueItems } from './key-value-list';

describe('stringArrayRecordToKeyValueItems', () => {
	it('maps a record of string arrays into joined key value items', () => {
		expect(
			stringArrayRecordToKeyValueItems({
				accept: ['application/json', 'text/plain'],
				encoding: ['gzip']
			})
		).toEqual([
			{ label: 'accept', value: 'application/json, text/plain' },
			{ label: 'encoding', value: 'gzip' }
		]);
	});
});
