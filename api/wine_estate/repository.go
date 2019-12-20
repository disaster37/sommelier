package wineEstate

import (
	"context"

	"github.com/sommelier/sommelier/models"
)

// Repository represent the address 's repository contract
type Repository interface {
	Fetch(ctx context.Context) ([]*models.WineEstate, error)
	GetByID(ctx context.Context, id string) (*models.WineEstate, error)
	Create(ctx context.Context, a *models.WineEstate) (string, error)
	Update(ctx context.Context, a *models.WineEstate) error
	Delete(ctx context.Context, id string) error
}
