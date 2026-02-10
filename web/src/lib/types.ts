export interface Source {
	ID: number;
	IngressUrl: string;
	EgressUrl: string;
	StaticHeaders: Record<string, string>;
	Status: string;
	StatusReason: string;
	Description: string;
	CreatedAt: Date;
	UpdatedAt: Date;
	DisableAt: Date | null;
}
