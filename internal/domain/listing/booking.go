package listing

import (
	"context"
	"time"
)

type BookingStatus string

const (
	Created   BookingStatus = "created"
	Confirmed BookingStatus = "confirmed"
	Cancelled BookingStatus = "cancelled"
	Completed BookingStatus = "completed"
)

type Booking struct {
	ID string
	// ID пользователя, создавшего бронирование (гость)
	UserID string
	// ID объявления, к которому относится бронирование
	ListingID string
	// Время бронирования
	CheckIn  time.Time
	Checkout time.Time
	Status   BookingStatus

	TotalPrice int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type BookingRepository interface {
	FindByID(ctx context.Context, id string) (*Booking, error)
	FindByUserID(ctx context.Context, userID string, limit, offset uint64) ([]*Booking, error)
	FindByStatus(ctx context.Context, status BookingStatus, limit, offset uint64) ([]*Booking, error)
	FindOverlapping(ctx context.Context, listingID string, checkIn, checkOut string) ([]*Booking, error)
	Save(ctx context.Context, booking *Booking) error
	DeleteByID(ctx context.Context, id string) error
	Count(ctx context.Context, status *BookingStatus) (int64, error)
}
