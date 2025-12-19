package user

import "context"

// Service defines domain-level user operations
type Service interface {
	// Authentication & Authorization
	Authenticate(ctx context.Context, email, password string) (*User, string, error) // returns user and token
	ValidateToken(ctx context.Context, token string) (*User, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)

	// User Management
	Create(ctx context.Context, email, phone, password string, role Role) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, id string, email, phone *string) (*User, error)
	Delete(ctx context.Context, id string) error

	// Profile Operations
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error

	// Email Verification
	SendVerificationEmail(ctx context.Context, userID string) error
	VerifyEmail(ctx context.Context, token string) error

	// Role Management
	SwitchRole(ctx context.Context, userID string, newRole Role) error
	RequestHostStatus(ctx context.Context, userID string) error
	ApproveHostStatus(ctx context.Context, userID string, approverID string) error

	// Admin Operations
	BanUser(ctx context.Context, userID string, reason string) error
	UnbanUser(ctx context.Context, userID string) error
	ListUsers(ctx context.Context, role *Role, limit, offset uint64) ([]*User, error)

	// Context Operations
	GetCurrentUser(ctx context.Context) (*User, error)
	GetUserStats(ctx context.Context, userID string) (*UserStats, error)
}

// UserStats represents user statistics
type UserStats struct {
	TotalBookings  int64
	TotalListings  int64
	ActiveBookings int64
	CompletedTrips int64
	Rating         float64
	MemberSince    string
}
