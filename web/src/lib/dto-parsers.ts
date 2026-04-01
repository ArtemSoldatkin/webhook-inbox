import type { DeliveryAttemptDTO, EventDTO, SourceDTO } from './types';

/** Raw DTO types with date fields as strings, representing the shape of API responses before parsing. */
type SourceDTORaw = Omit<SourceDTO, 'created_at' | 'updated_at' | 'disable_at'> & {
	created_at: string | Date;
	updated_at: string | Date;
	disable_at?: string | Date | null;
};

/** Converts a raw source payload into the app's typed source model. */
type EventDTORaw = Omit<EventDTO, 'received_at'> & {
	received_at: string | Date;
};

/** Raw DTO type for delivery attempts, with date fields as strings. */
type DeliveryAttemptDTORaw = Omit<
	DeliveryAttemptDTO,
	'started_at' | 'finished_at' | 'created_at' | 'next_attempt_at'
> & {
	started_at?: string | Date | null;
	finished_at?: string | Date | null;
	created_at: string | Date;
	next_attempt_at?: string | Date | null;
};

/**
 * Converts a raw source payload into the app's typed source model.
 *
 * @param raw - Source payload returned by the API.
 * @returns A source with date fields parsed.
 */
export function parseSourceDTO(raw: unknown): SourceDTO {
	const source = raw as SourceDTORaw;

	return {
		...source,
		created_at: new Date(source.created_at),
		updated_at: new Date(source.updated_at),
		disable_at: source.disable_at ? new Date(source.disable_at) : undefined
	};
}

/**
 * Converts a raw event payload into the app's typed event model.
 *
 * @param raw - Event payload returned by the API.
 * @returns An event with parsed timestamps.
 */
export function parseEventDTO(raw: unknown): EventDTO {
	const event = raw as EventDTORaw;

	return {
		...event,
		received_at: new Date(event.received_at)
	};
}

/**
 * Converts a raw delivery attempt payload into the typed UI model.
 *
 * @param raw - Delivery attempt payload returned by the API.
 * @returns A delivery attempt with parsed timestamps.
 */
export function parseDeliveryAttemptDTO(raw: unknown): DeliveryAttemptDTO {
	const attempt = raw as DeliveryAttemptDTORaw;

	return {
		...attempt,
		started_at: attempt.started_at ? new Date(attempt.started_at) : undefined,
		finished_at: attempt.finished_at ? new Date(attempt.finished_at) : undefined,
		created_at: new Date(attempt.created_at),
		next_attempt_at: attempt.next_attempt_at ? new Date(attempt.next_attempt_at) : undefined
	};
}
