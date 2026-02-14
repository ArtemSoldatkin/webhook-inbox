export type ContentType = 'application/json' | 'application/x-www-form-urlencoded' | 'text/plain';

export interface SourceDTO {
	ID: number;
	PublicID: string;
	IngressUrl: string;
	EgressUrl: string;
	StaticHeaders: Record<string, string>;
	Status: string;
	StatusReason: string;
	Description: string;
	CreatedAt: Date;
	UpdatedAt: Date;
	DisableAt: Date | null;
}

export interface EventDTO {
	ID: number;
	SourceID: number;
	Method: string;
	QueryParams: Record<string, string[]>;
	RawHeaders: Record<string, string[]>;
	Body: Record<string, string>;
	BodyContentType: string;
	ReceivedAt: Date;
}
