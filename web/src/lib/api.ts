import type { PaginatedResponse } from './types';

/**
 * Fetches a paginated API resource and normalizes the response shape.
 *
 * @param url - API endpoint without query parameters.
 * @param pageSize - Number of items to request per page.
 * @param nextCursor - Cursor for the next page, if present.
 * @param urlSearchParams - Additional filters appended to the request.
 * @returns The parsed paginated payload.
 */
export async function fetchPaginatedData<T>(
	url: string,
	pageSize: number,
	nextCursor: string | null,
	urlSearchParams: Record<string, string> = {}
): Promise<PaginatedResponse<T>> {
	const params = new URLSearchParams({
		limit: pageSize.toString(),
		...urlSearchParams
	});
	if (nextCursor) {
		params.append('cursor', nextCursor);
	}

	const response = await fetch(`${url}?${params.toString()}`);
	if (!response.ok) {
		throw new Error(`Failed to fetch data: ${response.statusText}`);
	}

	const result = await response.json();
	return {
		data: result.data as T[],
		next_cursor: result.next_cursor ?? null,
		has_next: result.has_next
	};
}
