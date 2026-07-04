package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
)

type Service interface {
	CreateProduct(ctx context.Context, params repository.CreateProductParams) (sql.Result, error)
	ListProducts(ctx context.Context, arg repository.ListProductsParams) (products []repository.Product, err error)
	ListProduct(ctx context.Context, id int64) (product repository.Product, err error)
	DeleteProduct(ctx context.Context, id int64) error
	UpdateProduct(ctx context.Context, arg repository.UpdateProductParams) (sql.Result, error)
}

type Svc struct {
	repository repository.Queries
}

type ListProductsFilter struct {
	Limit    int
	Offset   int
	Name     string
	MinPrice int
	MaxPrice int
}

func NewService(repository repository.Queries) Service {
	return &Svc{
		repository: repository,
	}
}

var (
	productNotFoundError = errors.New("Product not found")
	OneFieldRequired     = errors.New("At least one field is required")
	MinError             = errors.New("Price and quantity should be 0 and above")
)

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

func (svc *Svc) ListProducts(ctx context.Context, arg repository.ListProductsParams) (products []repository.Product, err error) {
	return svc.repository.ListProducts(ctx, arg)
}

func (svc *Svc) ListProduct(ctx context.Context, id int64) (product repository.Product, err error) {
	return svc.repository.ListProduct(ctx, id)
}

func (svc *Svc) DeleteProduct(ctx context.Context, id int64) error {
	return svc.repository.DeleteProduct(ctx, id)
}

func (svc *Svc) UpdateProduct(ctx context.Context, arg repository.UpdateProductParams) (sql.Result, error) {
	//& Validate fields
	if err := authutils.UpdateProductValidation(arg); err != nil {
		if authutils.ValidationErrorCheck("required", err) {
			return nil, fmt.Errorf("Product id required")
		}

		if authutils.ValidationErrorCheck("min", err) {
			return nil, MinError
		}

		return nil, err
	}
	
	//& Find product By id
	product, err := svc.repository.ListProduct(ctx, arg.ProductID)
	if err != nil {
		return nil, productNotFoundError
	}

	//& Ensure either price or quanity is provided
	if arg.Price == 0 && arg.Quantity == 0 {
		return nil, OneFieldRequired
	}

	//& provide default price and quantity
	if arg.Price == 0 {
		arg.Price = product.Price
	}

	if arg.Quantity == 0 {
		arg.Quantity = product.Quantity
	}
	return svc.repository.UpdateProduct(ctx, arg)
}
