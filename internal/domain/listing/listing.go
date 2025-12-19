package listing

import (
	"context"
	"time"
)

type Status string

const (
	Draft     Status = "draft"
	Published Status = "Published"
	Blocked   Status = "blocked"
	Deleted   Status = "deleted"
)

type Listing struct {
	ID          string
	Title       string
	Description string
	PricePerDay int64
	MinStayDays int64
	MaxStayDays int64
	Status      Status
	// ID изображений, связанных с объявлением
	ImagesIDs []string
	Address   Address
	//
	PublisherID string
	HostID      string

	CreatedAt time.Time
	UpdatedAT time.Time
}

type ListingRepository interface {
	FindByID(ctx context.Context, id string) (*Listing, error)
	FindByHostID(ctx context.Context, hostID string, limit, offset uint64) ([]*Listing, error)
	FindByStatus(ctx context.Context, status Status, limit, offset uint64) ([]*Listing, error)
	FindAll(ctx context.Context, limit, offset uint64) ([]*Listing, error)
	Save(ctx context.Context, listing *Listing) error
	DeleteByID(ctx context.Context, id string) error
	Count(ctx context.Context, status *Status) (int64, error)
	Exists(ctx context.Context, id string) (bool, error)
}
