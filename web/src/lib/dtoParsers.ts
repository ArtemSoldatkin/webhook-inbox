import type { DeliveryAttemptDTO, EventDTO, SourceDTO } from './types';

export function parseSourceDTO(raw: any): SourceDTO {
	return {
		...raw,
		created_at: new Date(raw.created_at),
		updated_at: new Date(raw.updated_at),
		disable_at: raw.disable_at && new Date(raw.disable_at)
	};
}

export function parseEventDTO(raw: any): EventDTO {
	return {
		...raw,
		received_at: new Date(raw.received_at)
	};
}

export function parseDeliveryAttemptDTO(raw: any): DeliveryAttemptDTO {
	return {
		...raw,
		started_at: raw.started_at && new Date(raw.started_at),
		finished_at: raw.finished_at && new Date(raw.finished_at),
		created_at: new Date(raw.created_at),
		next_attempt_at: raw.next_attempt_at && new Date(raw.next_attempt_at)
	};
}
