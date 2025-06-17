package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Asif-Faizal/Gommerce/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, userId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(userId),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyJWT verifies a JWT token and returns the user ID if valid
func VerifyJWT(tokenString string, secret []byte) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	// Check expiration
	expiredAt, ok := claims["expiredAt"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid expiration claim")
	}
	if float64(time.Now().Unix()) > expiredAt {
		return 0, fmt.Errorf("token has expired")
	}

	// Get user ID
	userIdStr, ok := claims["userId"].(string)
	if !ok {
		return 0, fmt.Errorf("invalid userId claim")
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return 0, fmt.Errorf("invalid userId format")
	}

	return userId, nil
}
