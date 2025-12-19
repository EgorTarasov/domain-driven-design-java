package user

import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/user"
)

// LoginUserInput represents the input for user login
type LoginUserInput struct {
	Email    string
	Password string
}

// LoginUserOutput represents the output of user login
type LoginUserOutput struct {
	User         *user.User
	AccessToken  string
	RefreshToken string
}

// LoginUserUseCase handles user login
type LoginUserUseCase struct {
	userRepo       user.Repository
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
}

// NewLoginUserUseCase creates a new LoginUserUseCase
func NewLoginUserUseCase(
	userRepo user.Repository,
	passwordHasher PasswordHasher,
	tokenGenerator TokenGenerator,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

// Execute authenticates a user
func (uc *LoginUserUseCase) Execute(ctx context.Context, input LoginUserInput) (*LoginUserOutput, error) {
	// Validate input
	if input.Email == "" || input.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	// Find user by email
	u, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	if err := uc.passwordHasher.Verify(input.Password, u.ID); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate tokens
	accessToken, err := uc.tokenGenerator.GenerateToken(u.ID, u.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.tokenGenerator.GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &LoginUserOutput{
		User:         u,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
