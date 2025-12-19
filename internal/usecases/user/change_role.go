package user

import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/user"
)

// ChangeRoleInput represents the input for changing user role
type ChangeRoleInput struct {
	UserID  string
	NewRole user.Role
}

// ChangeRoleUseCase handles changing a user's role
type ChangeRoleUseCase struct {
	userRepo user.Repository
}

// NewChangeRoleUseCase creates a new ChangeRoleUseCase
func NewChangeRoleUseCase(userRepo user.Repository) *ChangeRoleUseCase {
	return &ChangeRoleUseCase{
		userRepo: userRepo,
	}
}

// Execute changes a user's role
func (uc *ChangeRoleUseCase) Execute(ctx context.Context, input ChangeRoleInput) (*user.User, error) {
	// Validate input
	if input.UserID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Check authorization - only admin can change roles to admin
	currentRole, ok := ctx.Value("user_role").(user.Role)
	if !ok {
		return nil, fmt.Errorf("unauthorized: not authenticated")
	}

	if input.NewRole == user.Admin && currentRole != user.Admin {
		return nil, fmt.Errorf("unauthorized: only admins can promote users to admin")
	}

	// For switching between Guest and Host, user can do it themselves
	currentUserID, _ := ctx.Value("user_id").(string)
	if currentUserID != input.UserID && currentRole != user.Admin {
		return nil, fmt.Errorf("unauthorized: cannot change another user's role")
	}

	// Find user
	u, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Validate role transition
	if !isValidRoleTransition(u.Role, input.NewRole) {
		return nil, fmt.Errorf("invalid role transition from %s to %s", u.Role, input.NewRole)
	}

	// Update role
	u.Role = input.NewRole

	// Save updated user
	if err := uc.userRepo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}

	return u, nil
}

func isValidRoleTransition(from, to user.Role) bool {
	// Define valid role transitions
	validTransitions := map[user.Role][]user.Role{
		user.Guest: {user.Host},
		user.Host:  {user.Guest},
		user.Admin: {user.Host, user.Guest}, // Admin can downgrade
	}

	allowedRoles, ok := validTransitions[from]
	if !ok {
		return false
	}

	for _, role := range allowedRoles {
		if role == to {
			return true
		}
	}
	return false
}
