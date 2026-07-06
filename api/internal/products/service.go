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
	DeleteProduct(ctx context.Context, id int64) (repository.Product, error)
	UpdateProduct(ctx context.Context, arg repository.UpdateProductParams) (sql.Result, error)
}

type Svc struct {
	repository repository.Queries
	db         *sql.DB
}

type ListProductsFilter struct {
	Limit    int
	Offset   int
	Name     string
	MinPrice int
	MaxPrice int
}

func NewService(repository repository.Queries, db *sql.DB) Service {
	return &Svc{
		repository: repository,
		db:         db,
	}
}

var (
	productNotFoundError = errors.New("Product not found")
	OneFieldRequired     = errors.New("At least one field is required")
	MinError             = errors.New("Price and quantity should be 0 and above")
	DeleteAfterDelivery  = errors.New("Product still has an undelivered order")
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

// ! DELETE PRODUCT (ADMIN)
func (svc *Svc) DeleteProduct(ctx context.Context, id int64) (repository.Product, error) {
	// we need to ensure the product has already been paid for and marked "completed"
	// transaction for cancellinng the order and order_items for that product

	//& Ensure product exists
	product, err := svc.repository.ListProduct(ctx, id)
	if err != nil {
		return repository.Product{}, productNotFoundError
	}

	//& Ensure we don't delete product if there are paid orders(meaning they are paid for but not delivered)
	paid_orders, err := svc.repository.GetPaidProductOrders(ctx, id)
	if err != nil {
		return repository.Product{}, err
	}

	if len(paid_orders) > 0 {
		return repository.Product{}, DeleteAfterDelivery
	}

	//& Start transaction
	tx, err := svc.db.Begin()
	if err != nil {
		return repository.Product{}, err
	}

	defer tx.Rollback()

	qtx := svc.repository.WithTx(tx)

	//& Delete All order items with that product_id
	if err := qtx.DeleteOrderItemByProductID(ctx, id); err != nil {
		return repository.Product{}, err
	}

	//& Delete All orders with that product_id
	// - This is an issue because orders have no product id.
	// Find a way to combine order_items and orders and delete all records with that product id
	if err := qtx.DeleteOrderswithoutItems(ctx); err != nil {
		return repository.Product{}, err
	}

	//& Delete the product
	if err := qtx.DeleteProduct(ctx, id); err != nil {
		return repository.Product{}, err
	}

	if err := tx.Commit(); err != nil {
		return repository.Product{}, err
	}

	return product, nil
}

//! Update Product
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
