package user


import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/user"
)

// ListUsersInput represents the input for listing users
type ListUsersInput struct {
	Role   *user.Role
	Limit  uint64
	Offset uint64
}




















































}	}, nil		TotalCount: totalCount,		Users:      users,	return &ListUsersOutput{	}		return nil, fmt.Errorf("failed to count users: %w", err)	if err != nil {	totalCount, err := uc.userRepo.Count(ctx, input.Role)	// Get total count	}		return nil, fmt.Errorf("failed to retrieve users: %w", err)	if err != nil {	users, err := uc.userRepo.FindAll(ctx, input.Role, input.Limit, input.Offset)	// Get users	}		input.Limit = 100 // max limit	if input.Limit > 100 {	}		input.Limit = 20	if input.Limit == 0 {	// Set default pagination	}		return nil, fmt.Errorf("unauthorized: only admins can list users")	if !ok || currentRole != user.Admin {	currentRole, ok := ctx.Value("user_role").(user.Role)	// Check authorization - only admin can list usersfunc (uc *ListUsersUseCase) Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {// Execute retrieves a list of users with filters}	}		userRepo: userRepo,	return &ListUsersUseCase{func NewListUsersUseCase(userRepo user.Repository) *ListUsersUseCase {// NewListUsersUseCase creates a new ListUsersUseCase}	userRepo user.Repositorytype ListUsersUseCase struct {// ListUsersUseCase handles listing users with filters}	TotalCount int64	Users      []*user.Usertype ListUsersOutput struct {// ListUsersOutput represents the output of listing users