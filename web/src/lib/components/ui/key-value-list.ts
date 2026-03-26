/** Type definition for an individual key-value pair item. */
export type KeyValueItem = {
	/** The label or key for the item. This is a required string that represents the name or identifier of the item. */
	label: string;

	/** The value associated with the label. */
	value?: string | number | null;
};

/**
 * Utility function to convert a record of string keys and string array values into an array of key-value items.
 * Each key in the input record becomes a label, and the corresponding array of strings is joined into a single string value.
 *
 * @param data - A record where each key is a string and each value is an array of strings.
 * @returns An array of KeyValueItem objects, where each item has a label (the key) and a value (the joined string of the array).
 */
export function stringArrayRecordToKeyValueItems(data: Record<string, string[]>): KeyValueItem[] {
	return Object.entries(data).map(([label, value]) => ({
		label,
		value: value.join(', ')
	}));
}
