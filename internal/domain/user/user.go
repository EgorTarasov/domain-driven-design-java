package user

import "time"

type Role string

const (
	Guest Role = "guest"
	Host  Role = "host"
	Admin Role = "admin"
)

type User struct {
	ID        string
	Email     string
	Phone     string
	Role      Role
	CreatedAt time.Time
}
