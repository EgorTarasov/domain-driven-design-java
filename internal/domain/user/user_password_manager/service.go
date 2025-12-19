package userpasswordmanager

import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/user"
)

// Service handles user password management operations
type Service struct {
	userRepo       user.Repository
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
	emailService   EmailService
}

// NewService creates a new user password manager service
func NewService(
	userRepo user.Repository,
	passwordHasher PasswordHasher,
	tokenGenerator TokenGenerator,
	emailService EmailService,
) *Service {
	return &Service{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
		emailService:   emailService,
	}
}

// ChangePasswordInput represents the input for changing password
type ChangePasswordInput struct {
	UserID      string
	OldPassword string
	NewPassword string
}

// ChangePassword updates a user's password
func (s *Service) ChangePassword(ctx context.Context, input ChangePasswordInput) error {
	// Validate input
	if input.UserID == "" {
		return fmt.Errorf("user ID is required")
	}

	if input.OldPassword == "" || input.NewPassword == "" {
		return fmt.Errorf("old and new passwords are required")
	}

	if len(input.NewPassword) < 8 {
		return fmt.Errorf("new password must be at least 8 characters")
	}

	// Check authorization - user can only change their own password
	currentUserID, ok := ctx.Value("user_id").(string)
	if !ok || currentUserID != input.UserID {
		return fmt.Errorf("unauthorized: cannot change another user's password")
	}

	// Verify user exists
	exists, err := s.userRepo.Exists(ctx, input.UserID)
	if err != nil || !exists {
		return fmt.Errorf("user not found")
	}

	// TODO: Verify old password from password storage
	// if err := s.passwordHasher.Verify(input.OldPassword, storedHash); err != nil {
	// 	return fmt.Errorf("invalid old password")
	// }

	// Hash new password
	newHash, err := s.passwordHasher.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// TODO: Save new password hash to password storage
	_ = newHash

	return nil
}

// RequestPasswordReset sends a password reset email
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	// Validate input
	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Find user by email
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists
		return nil
	}

	// Generate reset token
	resetToken, err := s.tokenGenerator.GeneratePasswordResetToken(u.ID)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Send reset email
	if err := s.emailService.SendPasswordResetEmail(ctx, email, resetToken); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPasswordInput represents the input for resetting password
type ResetPasswordInput struct {
	Token       string
	NewPassword string
}

// ResetPassword resets a user's password using a reset token
func (s *Service) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	// Validate input
	if input.Token == "" || input.NewPassword == "" {
		return fmt.Errorf("token and new password are required")
	}

	if len(input.NewPassword) < 8 {
		return fmt.Errorf("new password must be at least 8 characters")
	}

	// Validate reset token
	userID, _, err := s.tokenGenerator.ValidateToken(input.Token)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Verify user exists
	exists, err := s.userRepo.Exists(ctx, userID)
	if err != nil || !exists {
		return fmt.Errorf("user not found")
	}

	// Hash new password
	newHash, err := s.passwordHasher.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// TODO: Save new password hash to password storage
	_ = newHash

	return nil
}

// PasswordHasher defines password hashing operations
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hashedPassword string) error
}

// TokenGenerator defines token generation operations
type TokenGenerator interface {
	GeneratePasswordResetToken(userID string) (string, error)
	ValidateToken(token string) (userID string, role user.Role, err error)
}

// EmailService defines email operations
type EmailService interface {
	SendPasswordResetEmail(ctx context.Context, email, resetToken string) error
}
