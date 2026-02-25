export type ContentType = 'application/json' | 'application/x-www-form-urlencoded' | 'text/plain';

export interface SourceDTO {
	ID: number;
	PublicID: string;
	IngressUrl: string;
	EgressUrl: string;
	StaticHeaders: Record<string, string>;
	Status: string;
	StatusReason: string | null;
	Description: string | null;
	CreatedAt: Date;
	UpdatedAt: Date;
	DisableAt: Date | null;
}

export interface EventDTO {
	ID: number;
	SourceID: number;
	Method: string;
	IngressPath: string;
	RemoteAddress: string | null;
	QueryParams: Record<string, string[]>;
	RawHeaders: Record<string, string[]>;
	Body: string; // Base64-encoded string
	BodyContentType: string;
	ReceivedAt: Date;
}

export interface DeliveryAttemptDTO {
	ID: number;
	EventID: number;
	AttemptNumber: number;
	State: string;
	StatusCode: number | null;
	ErrorType: string | null;
	ErrorMessage: string | null;
	StartedAt: Date | null;
	FinishedAt: Date | null;
	CreatedAt: Date;
	NextAttemptAt: Date | null;
}
