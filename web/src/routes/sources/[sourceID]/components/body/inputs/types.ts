/** Supported field types for multipart form construction. */
export type InputType = 'text' | 'number' | 'checkbox' | 'file' | 'date';

/** Shared shape for all dynamic form fields. */
interface BaseField {
	type: InputType;
	name: string;
}

/** Text-like field backed by a string value. */
interface TextField extends BaseField {
	type: 'text' | 'number' | 'date';
	value?: string | null;
}

/** Checkbox field backed by a boolean value. */
interface BooleanField extends BaseField {
	type: 'checkbox';
	value?: boolean;
}

/** File upload field backed by a file list. */
interface FileField extends BaseField {
	type: 'file';
	value?: FileList | null;
}

/** Any supported form field used by the constructor UI. */
export type FormField = TextField | BooleanField | FileField;
