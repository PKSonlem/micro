package entity

import "time"

type House struct {
	ID        int
	Address   string
	Year      int
	Developer *string
	CreatedAt time.Time
	UpdatedAt time.Time // last time flat added
}
