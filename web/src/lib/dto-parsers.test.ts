import { describe, expect, it } from 'vitest';
import { parseDeliveryAttemptDTO, parseEventDTO, parseSourceDTO } from './dto-parsers';

describe('parseSourceDTO', () => {
	it('parses source timestamps into Date instances', () => {
		const result = parseSourceDTO({
			id: 1,
			public_id: 'src_123',
			ingress_url: 'http://localhost/in',
			egress_url: 'http://localhost/out',
			status: 'active',
			created_at: '2026-04-01T10:00:00.000Z',
			updated_at: '2026-04-02T10:00:00.000Z',
			disable_at: '2026-04-03T10:00:00.000Z'
		});

		expect(result.created_at).toBeInstanceOf(Date);
		expect(result.updated_at).toBeInstanceOf(Date);
		expect(result.disable_at).toBeInstanceOf(Date);
		expect(result.disable_at?.toISOString()).toBe('2026-04-03T10:00:00.000Z');
	});

	it('preserves falsy optional disable_at values', () => {
		const result = parseSourceDTO({
			id: 1,
			public_id: 'src_123',
			ingress_url: 'http://localhost/in',
			egress_url: 'http://localhost/out',
			status: 'active',
			created_at: '2026-04-01T10:00:00.000Z',
			updated_at: '2026-04-02T10:00:00.000Z',
			disable_at: null
		});

		expect(result.disable_at).toBeNull();
	});
});

describe('parseEventDTO', () => {
	it('parses received_at into a Date instance', () => {
		const result = parseEventDTO({
			id: 11,
			source_id: 4,
			method: 'POST',
			ingress_path: '/hook',
			received_at: '2026-04-01T12:00:00.000Z'
		});

		expect(result.received_at).toBeInstanceOf(Date);
		expect(result.received_at.toISOString()).toBe('2026-04-01T12:00:00.000Z');
	});
});

describe('parseDeliveryAttemptDTO', () => {
	it('parses delivery attempt timestamps into Date instances', () => {
		const result = parseDeliveryAttemptDTO({
			id: 21,
			event_id: 11,
			attempt_number: 2,
			state: 'failed',
			started_at: '2026-04-01T12:00:00.000Z',
			finished_at: '2026-04-01T12:01:00.000Z',
			created_at: '2026-04-01T11:59:00.000Z',
			next_attempt_at: '2026-04-01T12:05:00.000Z'
		});

		expect(result.started_at).toBeInstanceOf(Date);
		expect(result.finished_at).toBeInstanceOf(Date);
		expect(result.created_at).toBeInstanceOf(Date);
		expect(result.next_attempt_at).toBeInstanceOf(Date);
	});

	it('preserves missing optional timestamps', () => {
		const result = parseDeliveryAttemptDTO({
			id: 21,
			event_id: 11,
			attempt_number: 2,
			state: 'pending',
			started_at: null,
			finished_at: null,
			created_at: '2026-04-01T11:59:00.000Z',
			next_attempt_at: null
		});

		expect(result.started_at).toBeNull();
		expect(result.finished_at).toBeNull();
		expect(result.next_attempt_at).toBeNull();
	});
});
