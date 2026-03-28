// web/src/lib/utils/cx.ts

import type { ClassValue } from 'svelte/elements';
import { twMerge } from 'tailwind-merge';

/**
 * Merges and deduplicates Tailwind CSS classes.
 * Usage: cx('p-2', condition && 'bg-red-500', 'text-sm')
 *
 * @param classNames - An array of class names, which can include falsy values.
 * @returns A string of merged and deduplicated class names.
 */
export function cx(...classNames: (ClassValue | string | false | null | undefined)[]) {
	return twMerge(classNames.filter(Boolean).join(' '));
}
