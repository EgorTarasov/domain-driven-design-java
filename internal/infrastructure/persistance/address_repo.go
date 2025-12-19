package persistance

import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/listing"

	sq "github.com/Masterminds/squirrel"
)

// AddressRepository handles address persistence operations
type AddressRepository struct {
	*MinimalRepository
}

var _ listing.AddressRepository = (*AddressRepository)(nil)

// NewAddressRepository creates a new address repository
func NewAddressRepository(db DB) *AddressRepository {
	return &AddressRepository{
		MinimalRepository: NewMinimalRepository(db),
	}
}

// FindByID finds an address by ID
func (r *AddressRepository) FindByID(ctx context.Context, id string) (*listing.Address, error) {
	query := r.GetBuilder().
		Select("id", "country", "city", "street", "house", "latitude", "longitude").
		From("addresses").
		Where(sq.Eq{"id": id})

	var a listing.Address
	err := r.QueryRow(ctx, query).Scan(
		&a.ID, &a.Country, &a.City, &a.Street, &a.House, &a.Latitude, &a.Longitude,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find address: %w", err)
	}

	return &a, nil
}

// Save creates or updates an address
func (r *AddressRepository) Save(ctx context.Context, a *listing.Address) error {
	query := r.GetBuilder().
		Insert("addresses").
		Columns("id", "country", "city", "street", "house", "latitude", "longitude").
		Values(a.ID, a.Country, a.City, a.Street, a.House, a.Latitude, a.Longitude).
		Suffix(`ON CONFLICT (id) DO UPDATE SET 
			country = EXCLUDED.country, 
			city = EXCLUDED.city, 
			street = EXCLUDED.street,
			house = EXCLUDED.house,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude`)

	_, err := r.Insert(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to save address: %w", err)
	}

	return nil
}

// DeleteByID deletes an address by ID
func (r *AddressRepository) DeleteByID(ctx context.Context, id string) error {
	query := r.GetBuilder().
		Delete("addresses").
		Where(sq.Eq{"id": id})

	result, err := r.Delete(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("address not found")
	}

	return nil
}
