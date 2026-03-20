export type InputType = 'text' | 'number' | 'checkbox' | 'file' | 'date';

interface BaseField {
	type: InputType;
	name: string;
}

interface TextField extends BaseField {
	type: 'text' | 'number' | 'date';
	value?: string | null;
}

interface BooleanField extends BaseField {
	type: 'checkbox';
	value?: boolean;
}

interface FileField extends BaseField {
	type: 'file';
	value?: FileList | null;
}

export type FormField = TextField | BooleanField | FileField;
