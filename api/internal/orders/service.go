package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
)

type Service interface {
	CreateOrder(ctx context.Context, params CreateOrderParams) (repository.Order, error)
	GetMyOrder(ctx context.Context) ([]repository.IdUserOrderDetailsRow, error)
	GetAllOrders(ctx context.Context) ([]repository.AllUsersOrderDetailsRow, error)
	CancelOrder(ctx context.Context, orderID int64) (sql.Result, error)
}

type Svc struct {
	repository repository.Queries
	db         *sql.DB
}

var (
	NotFoundError  = errors.New("Product Not Found")
	ProductNoStock = errors.New("The quantity selected is more than the stock available")
	AuthNotFound   = errors.New("Authentication context missing")
)

func NewService(repository repository.Queries, db *sql.DB) Service {
	return &Svc{
		repository: repository,
		db:         db,
	}
}

// ! CREATE ORDER
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

// ! GET MY ORDER
func (svc *Svc) GetMyOrder(ctx context.Context) ([]repository.IdUserOrderDetailsRow, error) {
	//& Ensure authenticated user exists from context
	user_id, ok := authutils.GetUserIDFromContext(ctx)
	if !ok {
		return []repository.IdUserOrderDetailsRow{}, AuthNotFound
	}
	return svc.repository.IdUserOrderDetails(ctx, user_id)
}

// ! GET ALL ORDERS (ADMIN)
func (svc *Svc) GetAllOrders(ctx context.Context) ([]repository.AllUsersOrderDetailsRow, error) {
	return svc.repository.AllUsersOrderDetails(ctx)
}

// ! CANCEL ORDER
func (svc *Svc) CancelOrder(ctx context.Context, orderID int64) (sql.Result, error) {
	//& Ensure authenticated user exists from context
	user_id, ok := authutils.GetUserIDFromContext(ctx)
	if !ok {
		return nil, AuthNotFound
	}
	//& Get Order
	order, err := svc.repository.GetOrderById(ctx, repository.GetOrderByIdParams{
		UserID:  user_id,
		OrderID: orderID,
	})
	if err != nil {
		return nil, err
	}

	if order.OrderStatus == "cancelled" {
		return nil, fmt.Errorf("Order already cancelled")
	}

	if order.OrderStatus == "paid" || order.OrderStatus == "completed" {
		return nil, fmt.Errorf("Order cannot be cancelled at this stage")
	}

	return svc.repository.CancelOrder(ctx, orderID)
}
