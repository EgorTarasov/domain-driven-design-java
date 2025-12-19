package persistance

import (
	"context"
	"database/sql"
	"fmt"

	"domain-driven-design-java/internal/domain/listing"

	sq "github.com/Masterminds/squirrel"
)

// ListingRepository handles listing persistence operations
type ListingRepository struct {
	*MinimalRepository
}

var _ listing.ListingRepository = (*ListingRepository)(nil)

// NewListingRepository creates a new listing repository
func NewListingRepository(db DB) *ListingRepository {
	return &ListingRepository{
		MinimalRepository: NewMinimalRepository(db),
	}
}

// FindByID finds a listing by ID
func (r *ListingRepository) FindByID(ctx context.Context, id string) (*listing.Listing, error) {
	query := r.GetBuilder().
		Select("id", "title", "description", "price_per_day", "min_stay_days", "max_stay_days",
			"status", "publisher_id", "host_id", "created_at", "updated_at").
		From("listings").
		Where(sq.Eq{"id": id})

	var l listing.Listing
	err := r.QueryRow(ctx, query).Scan(
		&l.ID, &l.Title, &l.Description, &l.PricePerDay, &l.MinStayDays, &l.MaxStayDays,
		&l.Status, &l.PublisherID, &l.HostID, &l.CreatedAt, &l.UpdatedAT,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find listing: %w", err)
	}

	return &l, nil
}

// FindByHostID finds listings by host ID
func (r *ListingRepository) FindByHostID(ctx context.Context, hostID string, limit, offset uint64) ([]*listing.Listing, error) {
	query := r.GetBuilder().
		Select("id", "title", "description", "price_per_day", "min_stay_days", "max_stay_days",
			"status", "publisher_id", "host_id", "created_at", "updated_at").
		From("listings").
		Where(sq.Eq{"host_id": hostID}).
		OrderBy("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	return r.queryListings(ctx, query)
}

// FindByStatus finds listings by status
func (r *ListingRepository) FindByStatus(ctx context.Context, status listing.Status, limit, offset uint64) ([]*listing.Listing, error) {
	query := r.GetBuilder().
		Select("id", "title", "description", "price_per_day", "min_stay_days", "max_stay_days",
			"status", "publisher_id", "host_id", "created_at", "updated_at").
		From("listings").
		Where(sq.Eq{"status": status}).
		OrderBy("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	return r.queryListings(ctx, query)
}

// FindAll retrieves all listings with pagination
func (r *ListingRepository) FindAll(ctx context.Context, limit, offset uint64) ([]*listing.Listing, error) {
	query := r.GetBuilder().
		Select("id", "title", "description", "price_per_day", "min_stay_days", "max_stay_days",
			"status", "publisher_id", "host_id", "created_at", "updated_at").
		From("listings").
		OrderBy("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	return r.queryListings(ctx, query)
}

// Save creates or updates a listing
func (r *ListingRepository) Save(ctx context.Context, l *listing.Listing) error {
	query := r.GetBuilder().
		Insert("listings").
		Columns("id", "title", "description", "price_per_day", "min_stay_days", "max_stay_days",
			"status", "publisher_id", "host_id", "created_at", "updated_at").
		Values(l.ID, l.Title, l.Description, l.PricePerDay, l.MinStayDays, l.MaxStayDays,
			l.Status, l.PublisherID, l.HostID, l.CreatedAt, l.UpdatedAT).
		Suffix(`ON CONFLICT (id) DO UPDATE SET 
			title = EXCLUDED.title, 
			description = EXCLUDED.description, 
			price_per_day = EXCLUDED.price_per_day,
			min_stay_days = EXCLUDED.min_stay_days,
			max_stay_days = EXCLUDED.max_stay_days,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at`)

	_, err := r.Insert(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to save listing: %w", err)
	}

	return nil
}

// DeleteByID deletes a listing by ID
func (r *ListingRepository) DeleteByID(ctx context.Context, id string) error {
	query := r.GetBuilder().
		Delete("listings").
		Where(sq.Eq{"id": id})

	result, err := r.Delete(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("listing not found")
	}

	return nil
}

// Count returns the total number of listings
func (r *ListingRepository) Count(ctx context.Context, status *listing.Status) (int64, error) {
	query := r.GetBuilder().
		Select("COUNT(*)").
		From("listings")

	if status != nil {
		query = query.Where(sq.Eq{"status": *status})
	}

	var count int64
	err := r.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count listings: %w", err)
	}

	return count, nil
}

// Exists checks if a listing exists by ID
func (r *ListingRepository) Exists(ctx context.Context, id string) (bool, error) {
	query := r.GetBuilder().
		Select("1").
		From("listings").
		Where(sq.Eq{"id": id}).
		Limit(1)

	var exists int
	err := r.QueryRow(ctx, query).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check listing existence: %w", err)
	}

	return true, nil
}

// queryListings executes a query and returns listings
func (r *ListingRepository) queryListings(ctx context.Context, query sq.SelectBuilder) ([]*listing.Listing, error) {
	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query listings: %w", err)
	}
	defer rows.Close()

	var listings []*listing.Listing
	for rows.Next() {
		var l listing.Listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.PricePerDay, &l.MinStayDays,
			&l.MaxStayDays, &l.Status, &l.PublisherID, &l.HostID, &l.CreatedAt, &l.UpdatedAT); err != nil {
			return nil, fmt.Errorf("failed to scan listing: %w", err)
		}
		listings = append(listings, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating listings: %w", err)
	}

	return listings, nil
}
