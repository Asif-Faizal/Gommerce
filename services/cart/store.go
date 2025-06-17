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

func (s *Store) GetOrders(userID int) ([]types.Order, error) {
	// First get all orders for the user
	query := `
		SELECT 
			o.id, 
			o.userId, 
			o.total, 
			o.status, 
			o.address, 
			o.createdAt,
			oi.id as item_id, 
			oi.orderId, 
			oi.productId, 
			oi.quantity, 
			oi.price,
			p.id as product_id, 
			p.name as product_name, 
			p.description as product_description, 
			p.image as product_image, 
			p.price as product_price, 
			p.quantity as product_quantity, 
			p.createdAt as product_createdAt
		FROM orders o
		LEFT JOIN order_items oi ON o.id = oi.orderId
		LEFT JOIN products p ON oi.productId = p.id
		WHERE o.userId = ?
		ORDER BY o.createdAt DESC, o.id ASC
	`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to store orders by ID
	ordersMap := make(map[int]*types.Order)

	for rows.Next() {
		var order types.Order
		var orderItem types.OrderItem
		var product types.Product
		var itemID sql.NullInt64
		var productID sql.NullInt64
		var productName sql.NullString
		var productDesc sql.NullString
		var productImage sql.NullString
		var productPrice sql.NullFloat64
		var productQuantity sql.NullInt32
		var productCreatedAt sql.NullTime

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Total,
			&order.Status,
			&order.Address,
			&order.CreatedAt,
			&itemID,
			&orderItem.OrderID,
			&orderItem.ProductID,
			&orderItem.Quantity,
			&orderItem.Price,
			&productID,
			&productName,
			&productDesc,
			&productImage,
			&productPrice,
			&productQuantity,
			&productCreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get or create order in map
		existingOrder, exists := ordersMap[order.ID]
		if !exists {
			order.Items = []types.OrderItem{}
			ordersMap[order.ID] = &order
			existingOrder = &order
		}

		// If there's an order item, add it to the order
		if itemID.Valid {
			orderItem.ID = int(itemID.Int64)
			if productID.Valid {
				product.ID = int(productID.Int64)
				product.Name = productName.String
				product.Description = productDesc.String
				product.Image = productImage.String
				product.Price = productPrice.Float64
				product.Quantity = int(productQuantity.Int32)
				product.CreatedAt = productCreatedAt.Time
				orderItem.Product = &product
			}
			existingOrder.Items = append(existingOrder.Items, orderItem)
		}
	}

	// Convert map to slice
	orders := make([]types.Order, 0, len(ordersMap))
	for _, order := range ordersMap {
		orders = append(orders, *order)
	}

	return orders, nil
}
