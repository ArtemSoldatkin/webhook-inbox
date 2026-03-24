/** Supported request body content types. */
export type ContentType =
	| 'application/json'
	| 'application/x-www-form-urlencoded'
	| 'text/plain'
	| 'multipart/form-data'
	| 'application/xml'
	| 'application/octet-stream';

/** Source entity returned by the backend API. */
export interface SourceDTO {
	/** Unique identifier for the source. */
	id: number;

	/** Public-facing identifier used in URLs. */
	public_id: string;

	/** URL where the source receives webhook events. */
	ingress_url: string;

	/** URL where the source forwards events. */
	egress_url: string;

	/** Optional headers added to forwarded requests. */
	static_headers?: Record<string, string>;

	/** Current lifecycle status of the source. */
	status: string;

	/** Optional explanation for the current source status. */
	status_reason?: string;

	/** User-provided description of the source. */
	description?: string;

	/** Time when the source was created. */
	created_at: Date;

	/** Time when the source was last updated. */
	updated_at: Date;

	/** Time when the source becomes disabled, if scheduled. */
	disable_at?: Date;
}

/** Captured inbound webhook event. */
export interface EventDTO {
	/** Unique identifier for the event. */
	id: number;

	/** Identifier of the source that received the event. */
	source_id: number;

	/** Optional hash used to detect duplicate events. */
	dedup_hash?: string;

	/** HTTP method used for the incoming request. */
	method: string;

	/** Request path used when the event was received. */
	ingress_path: string;

	/** Remote client address recorded for the request. */
	remote_address?: string;

	/** Parsed query parameters from the incoming request. */
	query_params?: Record<string, string[]>;

	/** Raw request headers grouped by header name. */
	raw_headers?: Record<string, string[]>;

	/** Request body stored as a Base64-encoded string. */
	body?: string;

	/** Content type associated with the stored request body. */
	body_content_type?: ContentType;

	/** Time when the event was received. */
	received_at: Date;
}

/** Delivery attempt recorded for a webhook event. */
export interface DeliveryAttemptDTO {
	/** Unique identifier for the delivery attempt. */
	id: number;

	/** Identifier of the event being delivered. */
	event_id: number;

	/** Sequential number of this delivery attempt. */
	attempt_number: number;

	/** Current state of the delivery attempt. */
	state: string;

	/** HTTP status code returned by the destination, if any. */
	status_code?: number;

	/** High-level error category for a failed attempt. */
	error_type?: string;

	/** Detailed error message for a failed attempt. */
	error_message?: string;

	/** Time when the attempt started. */
	started_at?: Date;

	/** Time when the attempt finished. */
	finished_at?: Date;

	/** Time when the attempt record was created. */
	created_at: Date;

	/** Time scheduled for the next retry, if applicable. */
	next_attempt_at?: Date;
}

/** Generic cursor-based API response. */
export interface PaginatedResponse<T> {
	/** Items returned for the current page. */
	data: T[];

	/** Cursor used to request the next page. */
	next_cursor: string | null;

	/** Indicates whether another page is available. */
	has_next: boolean;
}
