// Package user contains the user-related database operations
package user

import (
	"database/sql"
	"fmt"

	"github.com/Asif-Faizal/Gommerce/types"
)

// Store represents the user data store
// It implements the types.UserStore interface
type Store struct {
	db *sql.DB // Database connection
}

// NewStore creates a new instance of the user Store
// Takes a database connection as a parameter
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetUserByEmail retrieves a user from the database by their email
// Returns the user if found, or an error if not found or if there's a database error
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	query := "SELECT id, firstName, lastName, email, password, createdAt FROM users WHERE email = ?"
	user := &types.User{}

	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user from the database by their ID
// Returns the user if found, or an error if not found or if there's a database error
func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

// CreateUser inserts a new user into the database
// Takes a user object and returns any potential error
func (s *Store) CreateUser(user *types.User) error {
	query := `
		INSERT INTO users (firstName, lastName, email, password, createdAt)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.CreatedAt)
	return err
}

// scanRowsIntoUser is a helper function that scans database rows into a User struct
// Returns the user and any potential error
func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	if !rows.Next() {
		return nil, nil
	}

	user := &types.User{}
	if err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}
	return user, nil
}
