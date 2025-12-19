package userrolemanager

import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/user"
)

// Service handles user role management operations
type Service struct {
	userRepo user.Repository
}

// NewService creates a new user role manager service
func NewService(userRepo user.Repository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

// ChangeRoleInput represents the input for changing user role
type ChangeRoleInput struct {
	UserID  string
	NewRole user.Role
}

// ChangeRole changes a user's role
func (s *Service) ChangeRole(ctx context.Context, input ChangeRoleInput) (*user.User, error) {
	// Validate input
	if input.UserID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Check authorization
	if err := s.checkRoleChangeAuthorization(ctx, input.UserID, input.NewRole); err != nil {
		return nil, err
	}

	// Find user
	u, err := s.userRepo.FindByID(ctx, input.UserID)
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
	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}

	return u, nil
}

// SwitchToHost switches a guest user to host role
func (s *Service) SwitchToHost(ctx context.Context, userID string) (*user.User, error) {
	return s.ChangeRole(ctx, ChangeRoleInput{
		UserID:  userID,
		NewRole: user.Host,
	})
}

// SwitchToGuest switches a host user to guest role
func (s *Service) SwitchToGuest(ctx context.Context, userID string) (*user.User, error) {
	return s.ChangeRole(ctx, ChangeRoleInput{
		UserID:  userID,
		NewRole: user.Guest,
	})
}

// PromoteToAdmin promotes a user to admin role (admin only)
func (s *Service) PromoteToAdmin(ctx context.Context, userID string) (*user.User, error) {
	// Only admins can promote to admin
	currentRole, ok := ctx.Value("user_role").(user.Role)
	if !ok || currentRole != user.Admin {
		return nil, fmt.Errorf("unauthorized: only admins can promote users to admin")
	}

	return s.ChangeRole(ctx, ChangeRoleInput{
		UserID:  userID,
		NewRole: user.Admin,
	})
}

// checkRoleChangeAuthorization checks if the current user can change roles
func (s *Service) checkRoleChangeAuthorization(ctx context.Context, targetUserID string, newRole user.Role) error {
	currentRole, ok := ctx.Value("user_role").(user.Role)
	if !ok {
		return fmt.Errorf("unauthorized: not authenticated")
	}

	// Only admin can promote to admin
	if newRole == user.Admin && currentRole != user.Admin {
		return fmt.Errorf("unauthorized: only admins can promote users to admin")
	}

	// For switching between Guest and Host, user can do it themselves
	currentUserID, _ := ctx.Value("user_id").(string)
	if currentUserID != targetUserID && currentRole != user.Admin {
		return fmt.Errorf("unauthorized: cannot change another user's role")
	}

	return nil
}

// isValidRoleTransition checks if a role transition is valid
func isValidRoleTransition(from, to user.Role) bool {
	// Define valid role transitions
	validTransitions := map[user.Role][]user.Role{
		user.Guest: {user.Host},
		user.Host:  {user.Guest, user.Admin},
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
