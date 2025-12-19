package usermanager

import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/user"
)

// Service handles user management operations
type Service struct {
	userRepo user.Repository
}

// NewService creates a new user manager service
func NewService(userRepo user.Repository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

// GetByID retrieves a user by ID
func (s *Service) GetByID(ctx context.Context, userID string) (*user.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return u, nil
}

// GetByEmail retrieves a user by email
func (s *Service) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return u, nil
}

// GetCurrentUser retrieves the authenticated user from context
func (s *Service) GetCurrentUser(ctx context.Context) (*user.User, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user not authenticated")
	}

	return s.GetByID(ctx, userID)
}

// UpdateInput represents the input for updating a user
type UpdateInput struct {
	UserID string
	Email  *string
	Phone  *string
}

// Update updates a user's information
func (s *Service) Update(ctx context.Context, input UpdateInput) (*user.User, error) {
	// Validate input
	if input.UserID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Find existing user
	u, err := s.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check authorization
	if err := s.checkAuthorization(ctx, input.UserID); err != nil {
		return nil, err
	}

	// Update fields
	if input.Email != nil && *input.Email != "" {
		// Check if email is already taken
		existingUser, _ := s.userRepo.FindByEmail(ctx, *input.Email)
		if existingUser != nil && existingUser.ID != input.UserID {
			return nil, fmt.Errorf("email already in use")
		}
		u.Email = *input.Email
	}

	if input.Phone != nil {
		u.Phone = *input.Phone
	}

	// Save updated user
	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return u, nil
}

// Delete deletes or deactivates a user
func (s *Service) Delete(ctx context.Context, userID string) error {
	// Validate input
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Check authorization
	if err := s.checkAuthorization(ctx, userID); err != nil {
		return err
	}

	// Check if user exists
	exists, err := s.userRepo.Exists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Delete user
	if err := s.userRepo.DeleteByID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListInput represents the input for listing users
type ListInput struct {
	Role   *user.Role
	Limit  uint64
	Offset uint64
}

// ListOutput represents the output of listing users
type ListOutput struct {
	Users      []*user.User
	TotalCount int64
}

// List retrieves a list of users with filters
func (s *Service) List(ctx context.Context, input ListInput) (*ListOutput, error) {
	// Check authorization - only admin can list users
	currentRole, ok := ctx.Value("user_role").(user.Role)
	if !ok || currentRole != user.Admin {
		return nil, fmt.Errorf("unauthorized: only admins can list users")
	}

	// Set default pagination
	if input.Limit == 0 {
		input.Limit = 20
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	// Get users
	users, err := s.userRepo.FindAll(ctx, input.Role, input.Limit, input.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	// Get total count
	totalCount, err := s.userRepo.Count(ctx, input.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	return &ListOutput{
		Users:      users,
		TotalCount: totalCount,
	}, nil
}

// checkAuthorization checks if the current user can perform the operation
func (s *Service) checkAuthorization(ctx context.Context, targetUserID string) error {
	currentUserID, ok := ctx.Value("user_id").(string)
	if !ok || currentUserID != targetUserID {
		currentRole, _ := ctx.Value("user_role").(user.Role)
		if currentRole != user.Admin {
			return fmt.Errorf("unauthorized: cannot modify another user's data")
		}
	}
	return nil
}
