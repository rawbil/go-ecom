package products

import (
	"context"
	"database/sql"
	"fmt"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
)

type Service interface {
	CreateProduct(ctx context.Context, params repository.CreateProductParams) (sql.Result, error)
	ListProducts(ctx context.Context) (products []repository.Product, err error)
	ListProduct(ctx context.Context, id int64) (product repository.Product, err error)
	DeleteProduct(ctx context.Context, id int64) error
}

type Svc struct {
	repository repository.Queries
}

func NewService(repository repository.Queries) Service {
	return &Svc{
		repository: repository,
	}
}

func (svc *Svc) CreateProduct(ctx context.Context, params repository.CreateProductParams) (result sql.Result, error error) {
	// ensure both name and price are provided
	if params.ProductName == "" || params.Price == 0 {
		return result, fmt.Errorf("Product name and price must be provided")
	}

	// ensure price is positive
	if params.Price <= 0 {
		return result, fmt.Errorf("Price should be greater than 0")
	}

	return svc.repository.CreateProduct(ctx, params)
}

func (svc *Svc) ListProducts(ctx context.Context) (products []repository.Product, err error) {
	return svc.repository.ListProducts(ctx)
}

func (svc *Svc) ListProduct(ctx context.Context, id int64) (product repository.Product, err error) {
	return svc.repository.ListProduct(ctx, id)
}

func (svc *Svc) DeleteProduct(ctx context.Context, id int64) error {
	return svc.repository.DeleteProduct(ctx, id)
}
