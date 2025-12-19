package persistance

import (
	"context"
	"database/sql"
	"fmt"

	"domain-driven-design-java/internal/domain/user"

	sq "github.com/Masterminds/squirrel"
)

// UserRepository handles user persistence operations
type UserRepository struct {
	*MinimalRepository
}

var _ user.Repository = (*UserRepository)(nil)

// NewUserRepository creates a new user repository
func NewUserRepository(db DB) *UserRepository {
	return &UserRepository{
		MinimalRepository: NewMinimalRepository(db),
	}
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	query := r.GetBuilder().
		Select("id", "email", "phone", "role", "created_at").
		From("users").
		Where(sq.Eq{"id": id})

	var u user.User
	err := r.QueryRow(ctx, query).Scan(
		&u.ID,
		&u.Email,
		&u.Phone,
		&u.Role,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &u, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := r.GetBuilder().
		Select("id", "email", "phone", "role", "created_at").
		From("users").
		Where(sq.Eq{"email": email})

	var u user.User
	err := r.QueryRow(ctx, query).Scan(
		&u.ID,
		&u.Email,
		&u.Phone,
		&u.Role,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &u, nil
}

// FindAll retrieves all users with optional filters
func (r *UserRepository) FindAll(ctx context.Context, role *user.Role, limit, offset uint64) ([]*user.User, error) {
	query := r.GetBuilder().
		Select("id", "email", "phone", "role", "created_at").
		From("users").
		OrderBy("created_at DESC")

	if role != nil {
		query = query.Where(sq.Eq{"role": *role})
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Phone, &u.Role, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// Save creates or updates a user
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	query := r.GetBuilder().
		Insert("users").
		Columns("id", "email", "phone", "role", "created_at").
		Values(u.ID, u.Email, u.Phone, u.Role, u.CreatedAt).
		Suffix("ON CONFLICT (id) DO UPDATE SET email = EXCLUDED.email, phone = EXCLUDED.phone, role = EXCLUDED.role")

	_, err := r.Insert(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// DeleteByID deletes a user by ID
func (r *UserRepository) DeleteByID(ctx context.Context, id string) error {
	query := r.GetBuilder().
		Delete("users").
		Where(sq.Eq{"id": id})

	result, err := r.Delete(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Count returns the total number of users
func (r *UserRepository) Count(ctx context.Context, role *user.Role) (int64, error) {
	query := r.GetBuilder().
		Select("COUNT(*)").
		From("users")

	if role != nil {
		query = query.Where(sq.Eq{"role": *role})
	}

	var count int64
	err := r.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// Exists checks if a user exists by ID
func (r *UserRepository) Exists(ctx context.Context, id string) (bool, error) {
	query := r.GetBuilder().
		Select("1").
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1)

	var exists int
	err := r.QueryRow(ctx, query).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return true, nil
}
