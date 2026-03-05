export type ContentType = 'application/json' | 'application/x-www-form-urlencoded' | 'text/plain';

export interface SourceDTO {
	id: number;
	public_id: string;
	ingress_url: string;
	egress_url: string;
	static_headers?: Record<string, string>;
	status: string;
	status_reason?: string;
	description?: string;
	created_at: Date;
	updated_at: Date;
	disable_at?: Date;
}

export interface EventDTO {
	id: number;
	source_id: number;
	dedup_hash?: string;
	method: string;
	ingress_path: string;
	remote_address?: string;
	query_params?: Record<string, string[]>;
	raw_headers?: Record<string, string[]>;
	body?: string; // Base64-encoded string
	body_content_type?: string;
	received_at: Date;
}

export interface DeliveryAttemptDTO {
	id: number;
	event_id: number;
	attempt_number: number;
	state: string;
	status_code?: number;
	error_type?: string;
	error_message?: string;
	started_at?: Date;
	finished_at?: Date;
	created_at: Date;
	next_attempt_at?: Date;
}

export interface PaginatedResponse<T> {
	data: T[];
	next_cursor: string | null;
	has_next: boolean;
}
