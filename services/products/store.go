package products

import (
	"database/sql"
	"time"

	"github.com/Asif-Faizal/Gommerce/types"
)

// Store represents the user data store
// It implements the types.ProductStore interface
type Store struct {
	db *sql.DB // Database connection
}

// NewStore creates a new instance of the user Store
// Takes a database connection as a parameter
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetProductByID retrieves a product from the database by its ID
func (s *Store) GetProductByID() ([]types.Product, error) {
	query := "SELECT * FROM products"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := []types.Product{}
	for rows.Next() {
		product, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}
	return products, nil
}

// CreateProduct creates a new product in the database
func (s *Store) CreateProduct(product *types.Product) error {
	query := `
		INSERT INTO products (name, description, image, price, quantity, createdAt)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	// Set the creation time if not already set
	if product.CreatedAt.IsZero() {
		product.CreatedAt = time.Now()
	}

	result, err := s.db.Exec(
		query,
		product.Name,
		product.Description,
		product.Image,
		product.Price,
		product.Quantity,
		product.CreatedAt,
	)
	if err != nil {
		return err
	}

	// Get the ID of the newly created product
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = int(id)
	return nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := &types.Product{}
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return product, nil
}
