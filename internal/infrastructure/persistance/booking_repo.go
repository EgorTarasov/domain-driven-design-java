package persistance

import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/listing"

	sq "github.com/Masterminds/squirrel"
)

// BookingRepository handles booking persistence operations
type BookingRepository struct {
	*MinimalRepository
}

var _ listing.BookingRepository = (*BookingRepository)(nil)

// NewBookingRepository creates a new booking repository
func NewBookingRepository(db DB) *BookingRepository {
	return &BookingRepository{
		MinimalRepository: NewMinimalRepository(db),
	}
}

// FindByID finds a booking by ID
func (r *BookingRepository) FindByID(ctx context.Context, id string) (*listing.Booking, error) {
	query := r.GetBuilder().
		Select("id", "user_id", "listing_id", "check_in", "checkout", "status", "total_price", "created_at", "updated_at").
		From("bookings").
		Where(sq.Eq{"id": id})

	var b listing.Booking
	err := r.QueryRow(ctx, query).Scan(
		&b.ID, &b.UserID, &b.ListingID, &b.CheckIn, &b.Checkout, &b.Status, &b.TotalPrice, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find booking: %w", err)
	}

	return &b, nil
}

// FindByUserID finds bookings by user ID
func (r *BookingRepository) FindByUserID(ctx context.Context, userID string, limit, offset uint64) ([]*listing.Booking, error) {
	query := r.GetBuilder().
		Select("id", "user_id", "listing_id", "check_in", "checkout", "status", "total_price", "created_at", "updated_at").
		From("bookings").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	return r.queryBookings(ctx, query)
}

// FindByStatus finds bookings by status
func (r *BookingRepository) FindByStatus(ctx context.Context, status listing.BookingStatus, limit, offset uint64) ([]*listing.Booking, error) {
	query := r.GetBuilder().
		Select("id", "user_id", "listing_id", "check_in", "checkout", "status", "total_price", "created_at", "updated_at").
		From("bookings").
		Where(sq.Eq{"status": status}).
		OrderBy("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	return r.queryBookings(ctx, query)
}

// FindOverlapping finds bookings that overlap with the given date range
func (r *BookingRepository) FindOverlapping(ctx context.Context, listingID string, checkIn, checkOut string) ([]*listing.Booking, error) {
	query := r.GetBuilder().
		Select("id", "user_id", "listing_id", "check_in", "checkout", "status", "total_price", "created_at", "updated_at").
		From("bookings").
		Where(sq.Eq{"listing_id": listingID}).
		Where(sq.Or{
			sq.And{
				sq.LtOrEq{"check_in": checkOut},
				sq.GtOrEq{"checkout": checkIn},
			},
		}).
		Where(sq.NotEq{"status": listing.Cancelled})

	return r.queryBookings(ctx, query)
}

// Save creates or updates a booking
func (r *BookingRepository) Save(ctx context.Context, b *listing.Booking) error {
	query := r.GetBuilder().
		Insert("bookings").
		Columns("id", "user_id", "listing_id", "check_in", "checkout", "status", "total_price", "created_at", "updated_at").
		Values(b.ID, b.UserID, b.ListingID, b.CheckIn, b.Checkout, b.Status, b.TotalPrice, b.CreatedAt, b.UpdatedAt).
		Suffix(`ON CONFLICT (id) DO UPDATE SET 
			status = EXCLUDED.status,
			total_price = EXCLUDED.total_price,
			updated_at = EXCLUDED.updated_at`)

	_, err := r.Insert(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to save booking: %w", err)
	}

	return nil
}

// DeleteByID deletes a booking by ID
func (r *BookingRepository) DeleteByID(ctx context.Context, id string) error {
	query := r.GetBuilder().
		Delete("bookings").
		Where(sq.Eq{"id": id})

	result, err := r.Delete(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete booking: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

// Count returns the total number of bookings
func (r *BookingRepository) Count(ctx context.Context, status *listing.BookingStatus) (int64, error) {
	query := r.GetBuilder().
		Select("COUNT(*)").
		From("bookings")

	if status != nil {
		query = query.Where(sq.Eq{"status": *status})
	}

	var count int64
	err := r.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	return count, nil
}

// queryBookings executes a query and returns bookings
func (r *BookingRepository) queryBookings(ctx context.Context, query sq.SelectBuilder) ([]*listing.Booking, error) {
	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	var bookings []*listing.Booking
	for rows.Next() {
		var b listing.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.ListingID, &b.CheckIn, &b.Checkout, &b.Status,
			&b.TotalPrice, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bookings: %w", err)
	}

	return bookings, nil
}
