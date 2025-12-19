package persistance

import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/listing"

	sq "github.com/Masterminds/squirrel"
)

// AvailabilityRepository handles availability persistence operations
type AvailabilityRepository struct {
	*MinimalRepository
}

var _ listing.AvailabilityRepository = (*AvailabilityRepository)(nil)

// NewAvailabilityRepository creates a new availability repository
func NewAvailabilityRepository(db DB) *AvailabilityRepository {
	return &AvailabilityRepository{
		MinimalRepository: NewMinimalRepository(db),
	}
}

// FindByID finds an availability by ID
func (r *AvailabilityRepository) FindByID(ctx context.Context, id string) (*listing.Availability, error) {
	query := r.GetBuilder().
		Select("id", "date", "is_available", "price_override", "listing_id").
		From("availabilities").
		Where(sq.Eq{"id": id})

	var a listing.Availability
	err := r.QueryRow(ctx, query).Scan(&a.ID, &a.Date, &a.IsAvailable, &a.PriceOverride, &a.ListingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find availability: %w", err)
	}

	return &a, nil
}

// FindByListingIDAndDateRange finds availabilities for a listing within a date range
func (r *AvailabilityRepository) FindByListingIDAndDateRange(ctx context.Context, listingID string, startDate, endDate string) ([]*listing.Availability, error) {
	query := r.GetBuilder().
		Select("id", "date", "is_available", "price_override", "listing_id").
		From("availabilities").
		Where(sq.Eq{"listing_id": listingID}).
		Where(sq.GtOrEq{"date": startDate}).
		Where(sq.LtOrEq{"date": endDate}).
		OrderBy("date ASC")

	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query availabilities: %w", err)
	}
	defer rows.Close()

	var availabilities []*listing.Availability
	for rows.Next() {
		var a listing.Availability
		if err := rows.Scan(&a.ID, &a.Date, &a.IsAvailable, &a.PriceOverride, &a.ListingID); err != nil {
			return nil, fmt.Errorf("failed to scan availability: %w", err)
		}
		availabilities = append(availabilities, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating availabilities: %w", err)
	}

	return availabilities, nil
}

// Save creates or updates an availability
func (r *AvailabilityRepository) Save(ctx context.Context, a *listing.Availability) error {
	query := r.GetBuilder().
		Insert("availabilities").
		Columns("id", "date", "is_available", "price_override", "listing_id").
		Values(a.ID, a.Date, a.IsAvailable, a.PriceOverride, a.ListingID).
		Suffix(`ON CONFLICT (id) DO UPDATE SET 
			date = EXCLUDED.date,
			is_available = EXCLUDED.is_available,
			price_override = EXCLUDED.price_override,
			listing_id = EXCLUDED.listing_id`)

	_, err := r.Insert(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to save availability: %w", err)
	}

	return nil
}

// SaveMany creates or updates multiple availabilities
func (r *AvailabilityRepository) SaveMany(ctx context.Context, availabilities []*listing.Availability) error {
	if len(availabilities) == 0 {
		return nil
	}

	query := r.GetBuilder().
		Insert("availabilities").
		Columns("id", "date", "is_available", "price_override", "listing_id")

	for _, a := range availabilities {
		query = query.Values(a.ID, a.Date, a.IsAvailable, a.PriceOverride, a.ListingID)
	}

	query = query.Suffix(`ON CONFLICT (id) DO UPDATE SET 
		date = EXCLUDED.date,
		is_available = EXCLUDED.is_available,
		price_override = EXCLUDED.price_override,
		listing_id = EXCLUDED.listing_id`)

	_, err := r.Insert(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to save availabilities: %w", err)
	}

	return nil
}

// DeleteByID deletes an availability by ID
func (r *AvailabilityRepository) DeleteByID(ctx context.Context, id string) error {
	query := r.GetBuilder().
		Delete("availabilities").
		Where(sq.Eq{"id": id})

	result, err := r.Delete(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete availability: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("availability not found")
	}

	return nil

}
