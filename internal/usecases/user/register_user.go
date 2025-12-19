package user

import (
	"context"
	"fmt"
	"time"

	"domain-driven-design-java/internal/domain/user"

	"github.com/google/uuid"
)

// RegisterUserInput represents the input for user registration
type RegisterUserInput struct {
	Email    string
	Phone    string
	Password string
	Role     user.Role
}

// RegisterUserOutput represents the output of user registration
type RegisterUserOutput struct {
	User  *user.User
	Token string
}

// RegisterUserUseCase handles user registration
type RegisterUserUseCase struct {
	userRepo       user.Repository
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
	emailService   EmailService
}

// NewRegisterUserUseCase creates a new RegisterUserUseCase
func NewRegisterUserUseCase(
	userRepo user.Repository,
	passwordHasher PasswordHasher,
	tokenGenerator TokenGenerator,
	emailService EmailService,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
		emailService:   emailService,
	}
}

// Execute registers a new user
func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	// Validate input
	if err := uc.validateInput(input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Check if user already exists
	existingUser, _ := uc.userRepo.FindByEmail(ctx, input.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", input.Email)
	}

	// Hash password
	hashedPassword, err := uc.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	newUser := &user.User{
		ID:        uuid.New().String(),
		Email:     input.Email,
		Phone:     input.Phone,
		Role:      input.Role,
		CreatedAt: time.Now(),
	}

	// Save user and password
	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Generate authentication token
	token, err := uc.tokenGenerator.GenerateToken(newUser.ID, newUser.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Send verification email (async)
	go uc.emailService.SendVerificationEmail(context.Background(), newUser.Email, newUser.ID)

	return &RegisterUserOutput{
		User:  newUser,
		Token: token,
	}, nil
}

func (uc *RegisterUserUseCase) validateInput(input RegisterUserInput) error {
	if input.Email == "" {
		return fmt.Errorf("email is required")
	}
	if input.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(input.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if input.Role != user.Guest && input.Role != user.Host {
		return fmt.Errorf("invalid role")
	}
	return nil
}
