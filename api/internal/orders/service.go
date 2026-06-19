package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
)

type Service interface {
	CreateOrder(ctx context.Context, params CreateOrderParams) (repository.Order, error)
}

type Svc struct {
	repository repository.Queries
	db         *sql.DB
}

var (
	NotFoundError  = errors.New("Product Not Found")
	ProductNoStock = errors.New("The quantity selected is more than the stock available")
)

func NewService(repository repository.Queries, db *sql.DB) Service {
	return &Svc{
		repository: repository,
		db:         db,
	}
}

func (svc *Svc) CreateOrder(ctx context.Context, params CreateOrderParams) (result repository.Order, err error) {
	if params.UserID == 0 {
		return repository.Order{}, fmt.Errorf("User id not found")
	}

	if len(params.Items) < 1 {
		return repository.Order{}, fmt.Errorf("No order items found")
	}

	//? TRANSACTION
	tx, err := svc.db.Begin()
	if err != nil {
		return repository.Order{}, err
	}

	defer tx.Rollback()

	qtx := svc.repository.WithTx(tx)

	order, err := qtx.CreateOrder(ctx, params.UserID)
	if err != nil {
		return repository.Order{}, err
	}

	orderId, err := order.LastInsertId()
	if err != nil {
		return repository.Order{}, err
	}

	// ensure the products exist
	for _, item := range params.Items {
		product, err := qtx.ListProduct(ctx, item.ProductId)
		if err != nil {
			return repository.Order{}, NotFoundError
		}

		if product.Quantity < item.Quantity {
			return repository.Order{}, ProductNoStock
		}

		// Create order item
		_, err = qtx.CreateOrderItem(ctx, repository.CreateOrderItemParams{
			OrderID:    orderId,
			ProductID:  product.ProductID,
			Quantity:   item.Quantity,
			TotalPrice: int64(product.Price) * int64(item.Quantity),
		})

		if err != nil {
			return repository.Order{}, err
		}

		// reduce the quantity of the product

		updated_qty := product.Quantity - item.Quantity
		fmt.Println(updated_qty)

		if _, err = qtx.UpdateProductQuantity(ctx, repository.UpdateProductQuantityParams{
			Quantity:  updated_qty,
			ProductID: product.ProductID,
		}); err != nil {
			return repository.Order{}, err
		}

	}

	createdOrder, err := qtx.ListOrder(ctx, orderId)
	if err != nil {
		return repository.Order{}, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return repository.Order{}, err
	}

	return createdOrder, nil
}
