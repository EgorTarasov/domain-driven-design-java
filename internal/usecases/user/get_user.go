package user

import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/user"
)

// GetUserByIDUseCase handles retrieving a user by ID
type GetUserByIDUseCase struct {
	userRepo user.Repository
}

// NewGetUserByIDUseCase creates a new GetUserByIDUseCase
func NewGetUserByIDUseCase(userRepo user.Repository) *GetUserByIDUseCase {
	return &GetUserByIDUseCase{
		userRepo: userRepo,
	}
}

// Execute retrieves a user by ID
func (uc *GetUserByIDUseCase) Execute(ctx context.Context, userID string) (*user.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	u, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return u, nil
}

// GetCurrentUserUseCase handles retrieving the current authenticated user
type GetCurrentUserUseCase struct {
	userRepo user.Repository
}

// NewGetCurrentUserUseCase creates a new GetCurrentUserUseCase
func NewGetCurrentUserUseCase(userRepo user.Repository) *GetCurrentUserUseCase {
	return &GetCurrentUserUseCase{
		userRepo: userRepo,
	}
}

// Execute retrieves the current user from context
func (uc *GetCurrentUserUseCase) Execute(ctx context.Context) (*user.User, error) {
	// Extract user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user not authenticated")
	}

	u, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return u, nil
}
