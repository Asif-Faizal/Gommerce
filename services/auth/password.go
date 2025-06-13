// Package auth handles authentication-related functionality
package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword takes a plain text password and returns a bcrypt hashed version
// Uses bcrypt's default cost for hashing
// Returns the hashed password as a string and any potential error
func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash of the password using the default cost
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// Convert the byte slice to string before returning
	return string(hashedPassword), nil
}

// ComparePasswords compares a hashed password with a plain text password
// Returns true if the passwords match, false otherwise
func ComparePasswords(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
