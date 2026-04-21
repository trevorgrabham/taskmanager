package account

import "time"

type User struct {
	Username       string
	HashedPassword string
	ID             int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
