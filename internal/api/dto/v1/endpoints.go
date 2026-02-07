package dtov1

import "github.com/jackc/pgx/v5/pgtype"


type Endpoint struct {
	ID          int64
	UserID      pgtype.Int8
	Url         string
	Name        string
	Description pgtype.Text
	Headers     map[string]any
	IsActive    pgtype.Bool
	CreatedAt   pgtype.Timestamp
}
