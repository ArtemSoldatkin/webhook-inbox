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

export interface Webhook {
	ID: number;
	EndpointID: number;
	PublicKey: string;
	Name: string;
	Description: string;
	IsActive: boolean;
	CreatedAt: Date;
	UpdatedAt: Date;
}
