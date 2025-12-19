package usercreator

import (
	"context"
	"fmt"
	"time"

	"domain-driven-design-java/internal/domain/user"
)

// Service handles user creation operations
type Service struct {
	userRepo       user.Repository
	passwordHasher PasswordHasher
}

// NewService creates a new user creator service
func NewService(
	userRepo user.Repository,
	passwordHasher PasswordHasher,
) *Service {
	return &Service{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

// CreateInput represents the input for creating a user
type CreateInput struct {
	Email    string
	Password string
	Phone    string
	Role     user.Role
}

// Create creates a new user
func (s *Service) Create(ctx context.Context, input CreateInput) (*user.User, error) {
	// Validate input
	if input.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if input.Password == "" {
		return nil, fmt.Errorf("password is required")
	}
	if len(input.Password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters")
	}

	// Set default role if not provided
	if input.Role == "" {
		input.Role = user.Guest
	}

	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already in use")
	}

	// Hash password
	hashedPassword, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// TODO: Store password hash separately
	_ = hashedPassword

	// Create user
	u := &user.User{
		Email:     input.Email,
		Phone:     input.Phone,
		Role:      input.Role,
		CreatedAt: time.Now(),
	}

	// Save user
	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return u, nil
}

// PasswordHasher defines password hashing operations
type PasswordHasher interface {
	Hash(password string) (string, error)
}
