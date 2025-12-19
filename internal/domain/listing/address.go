package listing

import "context"

type Address struct {
	ID        string
	Country   string
	City      string
	Street    string
	House     string
	Latitude  float64
	Longitude float64
}

type AddressRepository interface {
	FindByID(ctx context.Context, id string) (*Address, error)
	Save(ctx context.Context, address *Address) error
	DeleteByID(ctx context.Context, id string) error
}
