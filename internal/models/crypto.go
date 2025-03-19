package models

import "time"

type Crypto struct {
	ID            int64     `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Balance       string    `json:"balance" db:"balance"`
	LastUpdatedAt time.Time `json:"last_updated_at" db:"last_updated_at"`
}
