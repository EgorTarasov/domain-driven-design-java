package user

import (
	"context"

	"domain-driven-design-java/internal/domain/user"
)

// PasswordHasher defines password hashing operations
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, userID string) error
}

// TokenGenerator defines token generation operations
type TokenGenerator interface {
	GenerateToken(userID string, role user.Role) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string) (userID string, role user.Role, err error)
}

// EmailService defines email operations
type EmailService interface {
	SendVerificationEmail(ctx context.Context, email, userID string) error
	SendPasswordResetEmail(ctx context.Context, email, token string) error
	SendWelcomeEmail(ctx context.Context, email, userName string) error
}
