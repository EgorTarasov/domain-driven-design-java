package persistance

import (
	"context"
	"fmt"

	"domain-driven-design-java/internal/domain/listing"

	sq "github.com/Masterminds/squirrel"
)

// ImageRepository handles image persistence operations
type ImageRepository struct {
	*MinimalRepository
}

var _ listing.ImageRepository = (*ImageRepository)(nil)

// NewImageRepository creates a new image repository
func NewImageRepository(db DB) *ImageRepository {
	return &ImageRepository{
		MinimalRepository: NewMinimalRepository(db),
	}
}

// FindByID finds an image by ID
func (r *ImageRepository) FindByID(ctx context.Context, id string) (*listing.Image, error) {
	query := r.GetBuilder().
		Select("id", "url", "position").
		From("images").
		Where(sq.Eq{"id": id})

	var img listing.Image
	err := r.QueryRow(ctx, query).Scan(&img.ID, &img.URL, &img.Position)
	if err != nil {
		return nil, fmt.Errorf("failed to find image: %w", err)
	}

	return &img, nil
}

// FindByIDs finds multiple images by their IDs
func (r *ImageRepository) FindByIDs(ctx context.Context, ids []string) ([]*listing.Image, error) {
	if len(ids) == 0 {
		return []*listing.Image{}, nil
	}

	query := r.GetBuilder().
		Select("id", "url", "position").
		From("images").
		Where(sq.Eq{"id": ids}).
		OrderBy("position ASC")

	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query images: %w", err)
	}
	defer rows.Close()

	var images []*listing.Image
	for rows.Next() {
		var img listing.Image
		if err := rows.Scan(&img.ID, &img.URL, &img.Position); err != nil {
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}
		images = append(images, &img)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating images: %w", err)
	}

	return images, nil
}

// SaveMany creates or updates multiple images
func (r *ImageRepository) SaveMany(ctx context.Context, images []*listing.Image) error {
	if len(images) == 0 {
		return nil
	}

	query := r.GetBuilder().
		Insert("images").
		Columns("id", "url", "position")

	for _, img := range images {
		query = query.Values(img.ID, img.URL, img.Position)
	}

	query = query.Suffix(`ON CONFLICT (id) DO UPDATE SET 
		url = EXCLUDED.url, 
		position = EXCLUDED.position`)

	_, err := r.Insert(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to save images: %w", err)
	}

	return nil
}

// DeleteByID deletes an image by ID
func (r *ImageRepository) DeleteByID(ctx context.Context, id string) error {
	query := r.GetBuilder().
		Delete("images").
		Where(sq.Eq{"id": id})

	result, err := r.Delete(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("image not found")
	}

	return nil
}

// DeleteByIDs deletes multiple images by their IDs
func (r *ImageRepository) DeleteByIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	query := r.GetBuilder().
		Delete("images").
		Where(sq.Eq{"id": ids})

	_, err := r.Delete(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete images: %w", err)
	}

	return nil
}
