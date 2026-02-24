import type { DeliveryAttemptDTO, EventDTO, SourceDTO } from './types';

export function parseSourceDTO(raw: any): SourceDTO {
	return {
		...raw,
		CreatedAt: new Date(raw.CreatedAt),
		UpdatedAt: new Date(raw.UpdatedAt),
		DisableAt: raw.DisableAt ? new Date(raw.DisableAt) : null
	};
}

export function parseEventDTO(raw: any): EventDTO {
	return {
		...raw,
		ReceivedAt: new Date(raw.ReceivedAt)
	};
}

export function parseDeliveryAttemptDTO(raw: any): DeliveryAttemptDTO {
	return {
		...raw,
		StartedAt: raw.StartedAt ? new Date(raw.StartedAt) : null,
		FinishedAt: raw.FinishedAt ? new Date(raw.FinishedAt) : null,
		CreatedAt: new Date(raw.CreatedAt),
		NextAttemptAt: raw.NextAttemptAt ? new Date(raw.NextAttemptAt) : null
	};
}
