export interface User {
	ID: number;
	Email: string;
	CreatedAt: Date;
}

export interface Endpoint {
	ID: number;
	UserID: number;
	Url: string;
	Name: string;
	Description: string;
	Headers: string[];
	IsActive: boolean;
	CreatedAt: Date;
}
