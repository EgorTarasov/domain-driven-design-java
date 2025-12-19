package userauthenticator

import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/user"
)

// Service handles user authentication operations
type Service struct {
	userRepo       user.Repository
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
}

// NewService creates a new user authenticator service
func NewService(
	userRepo user.Repository,
	passwordHasher PasswordHasher,
	tokenGenerator TokenGenerator,
) *Service {
	return &Service{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

// LoginInput represents the input for user login
type LoginInput struct {
	Email    string
	Password string
}

// LoginOutput represents the output of user login
type LoginOutput struct {
	User         *user.User
	AccessToken  string
	RefreshToken string
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// Validate input
	if input.Email == "" || input.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	// Find user by email
	u, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// TODO: Verify password from separate password storage
	// if err := s.passwordHasher.Verify(input.Password, storedHash); err != nil {
	// 	return nil, fmt.Errorf("invalid email or password")
	// }

	// Generate tokens
	accessToken, err := s.tokenGenerator.GenerateToken(u.ID, u.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenGenerator.GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &LoginOutput{
		User:         u,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateToken validates an authentication token
func (s *Service) ValidateToken(ctx context.Context, token string) (*user.User, error) {
	// Validate token
	userID, _, err := s.tokenGenerator.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return u, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Validate refresh token
	userID, role, err := s.tokenGenerator.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if user still exists
	exists, err := s.userRepo.Exists(ctx, userID)
	if err != nil || !exists {
		return "", fmt.Errorf("user not found")
	}

	// Generate new access token
	newToken, err := s.tokenGenerator.GenerateToken(userID, role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return newToken, nil
}

// PasswordHasher defines password hashing operations
type PasswordHasher interface {
	Verify(password, hashedPassword string) error
}

// TokenGenerator defines token generation operations
type TokenGenerator interface {
	GenerateToken(userID string, role user.Role) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string) (userID string, role user.Role, err error)
}