package listing

import (
	"context"
	"time"
)

type Availability struct {
	// UUID
	ID          string
	Date        time.Time
	IsAvailable bool
	// Если установлено, то цена на эту дату переопределена
	PriceOverride int64
	// ID объявления, к которому относится эта доступность
	ListingID string
}

type AvailabilityRepository interface {
	FindByID(ctx context.Context, id string) (*Availability, error)
	FindByListingIDAndDateRange(ctx context.Context, listingID string, startDate, endDate string) ([]*Availability, error)
	Save(ctx context.Context, availability *Availability) error
	SaveMany(ctx context.Context, availabilities []*Availability) error
	DeleteByID(ctx context.Context, id string) error
}
