package entity

import "time"

// TODO: возможно стоит выделить отдульную сущность CreateHouse, чтобы явно указать валидацию
type House struct {
	ID        int
	Address   string `validate:"required,min=5"`
	Year      int    `validate:"required,min=1900"`
	Developer *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
