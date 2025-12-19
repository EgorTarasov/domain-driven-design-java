package listing

import "context"

type Image struct {
	ID       string
	URL      string
	Position int64
}

type ImageRepository interface {
	FindByID(ctx context.Context, id string) (*Image, error)
	FindByIDs(ctx context.Context, ids []string) ([]*Image, error)
	SaveMany(ctx context.Context, images []*Image) error
	DeleteByID(ctx context.Context, id string) error
	DeleteByIDs(ctx context.Context, ids []string) error
}
