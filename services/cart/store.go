package cart

import (
	"database/sql"

	"github.com/Asif-Faizal/Gommerce/types"
)

// Store represents the user data store
// It implements the types.CartStore interface
type Store struct {
	db *sql.DB // Database connection
}

// NewStore creates a new instance of the user Store
// Takes a database connection as a parameter
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order *types.Order) (int, error) {
	query := "INSERT INTO orders (userId, total, status, address) VALUES (?, ?, ?, ?)"
	result, err := s.db.Exec(query, order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}
	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(orderID), nil
}

func (s *Store) CreateOrderItem(orderItem *types.OrderItem) error {
	query := "INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?, ?, ?, ?)"
	_, err := s.db.Exec(query, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	if err != nil {
		return err
	}
	return nil
}
