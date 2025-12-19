package user

import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/user"
)

// DeleteUserUseCase handles user deletion/deactivation
type DeleteUserUseCase struct {
	userRepo user.Repository
}

// NewDeleteUserUseCase creates a new DeleteUserUseCase
func NewDeleteUserUseCase(userRepo user.Repository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepo: userRepo,
	}
}

// Execute deletes/deactivates a user
func (uc *DeleteUserUseCase) Execute(ctx context.Context, userID string) error {
	// Validate input
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Check authorization - user can delete their own account or admin can delete any
	currentUserID, ok := ctx.Value("user_id").(string)
	if !ok || currentUserID != userID {
		currentRole, _ := ctx.Value("user_role").(user.Role)
		if currentRole != user.Admin {
			return fmt.Errorf("unauthorized: cannot delete another user's account")
		}
	}

	// Check if user exists
	exists, err := uc.userRepo.Exists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	// TODO: Check if user has active bookings or listings
	// This should be handled by domain logic or additional checks

	// Delete user
	if err := uc.userRepo.DeleteByID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
