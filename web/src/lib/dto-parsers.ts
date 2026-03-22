import type { DeliveryAttemptDTO, EventDTO, SourceDTO } from './types';

/**
 * Converts a raw source payload into the app's typed source model.
 *
 * @param raw - Source payload returned by the API.
 * @returns A source with date fields parsed.
 */
export function parseSourceDTO(raw: any): SourceDTO {
	return {
		...raw,
		created_at: new Date(raw.created_at),
		updated_at: new Date(raw.updated_at),
		disable_at: raw.disable_at && new Date(raw.disable_at)
	};
}

/**
 * Converts a raw event payload into the app's typed event model.
 *
 * @param raw - Event payload returned by the API.
 * @returns An event with parsed timestamps.
 */
export function parseEventDTO(raw: any): EventDTO {
	return {
		...raw,
		received_at: new Date(raw.received_at)
	};
}

/**
 * Converts a raw delivery attempt payload into the typed UI model.
 *
 * @param raw - Delivery attempt payload returned by the API.
 * @returns A delivery attempt with parsed timestamps.
 */
export function parseDeliveryAttemptDTO(raw: any): DeliveryAttemptDTO {
	return {
		...raw,
		started_at: raw.started_at && new Date(raw.started_at),
		finished_at: raw.finished_at && new Date(raw.finished_at),
		created_at: new Date(raw.created_at),
		next_attempt_at: raw.next_attempt_at && new Date(raw.next_attempt_at)
	};
}
