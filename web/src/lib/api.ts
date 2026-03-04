import type { PaginatedResponse } from './types';

export async function fetchPaginatedData<T>(
	url: string,
	pageSize: number,
	nextCursor: string | null
): Promise<PaginatedResponse<T>> {
	const params = new URLSearchParams({
		limit: pageSize.toString(),
		cursor: nextCursor || ''
	});
	const response = await fetch(`${url}?${params.toString()}`);
	if (!response.ok) {
		throw new Error(`Failed to fetch data: ${response.statusText}`);
	}
	const result = await response.json();
	return {
		data: result.data as T[],
		next_cursor: result.next_cursor,
		has_next: result.has_next ?? null
	};
}
