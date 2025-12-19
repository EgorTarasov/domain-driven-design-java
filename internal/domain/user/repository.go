package user

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindAll(ctx context.Context, role *Role, limit, offset uint64) ([]*User, error)
	Save(ctx context.Context, user *User) error
	DeleteByID(ctx context.Context, id string) error
	Count(ctx context.Context, role *Role) (int64, error)
	Exists(ctx context.Context, id string) (bool, error)
}
