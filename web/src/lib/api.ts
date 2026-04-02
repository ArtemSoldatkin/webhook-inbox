import type { PaginatedResponse } from './types';

/**
 * Generically formats a response's status and status text into a single string for error messages.
 *
 * @param response - The fetch Response object to summarize.
 * @returns A string summarizing the response status, e.g. "400 Bad Request", "500", "OK",
 * 	or a generic message if neither is available.
 */
function getResponseStatusSummary(response: Response): string {
	if (response.status && response.statusText) {
		return `${response.status} ${response.statusText}`;
	}
	if (response.status) {
		return String(response.status);
	}
	if (response.statusText) {
		return response.statusText;
	}
	return 'Request failed';
}

/**
 * Attempts to extract a meaningful error message from a failed API response, with fallbacks.
 *
 * @param response - The fetch Response object from the failed request.
 * @param fallbackPrefix - A prefix for the error message if the response body doesn't contain one.
 * @returns A string error message extracted from the response or a fallback message.
 */
export async function getResponseErrorMessage(
	response: Response,
	fallbackPrefix: string
): Promise<string> {
	let errorMsg = `${fallbackPrefix}: ${getResponseStatusSummary(response)}`;
	try {
		const errorJson = await response.json();
		if (errorJson && typeof errorJson.error === 'string' && errorJson.error.length > 0) {
			errorMsg = errorJson.error;
		}
	} catch {
		console.warn('Failed to parse error response as JSON', response);
	}

	return errorMsg;
}

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
		throw new Error(await getResponseErrorMessage(response, 'Failed to fetch data'));
	}

	const result = await response.json();
	return {
		data: result.data as T[],
		next_cursor: result.next_cursor ?? null,
		has_next: result.has_next
	};
}
