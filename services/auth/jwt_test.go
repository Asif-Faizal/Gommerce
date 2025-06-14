package auth

import (
	"testing"

	"github.com/Asif-Faizal/Gommerce/config"
	"github.com/golang-jwt/jwt"
)

func TestCreateJWT(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		secret   []byte
		userId   int
		wantErr  bool
		validate func(t *testing.T, token string)
	}{
		{
			name:    "successful token creation",
			secret:  []byte("test-secret"),
			userId:  123,
			wantErr: false,
			validate: func(t *testing.T, token string) {
				if token == "" {
					t.Error("expected non-empty token")
				}
				// Parse the token to verify its contents
				parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return []byte("test-secret"), nil
				})
				if err != nil {
					t.Errorf("failed to parse token: %v", err)
				}
				if !parsedToken.Valid {
					t.Error("expected valid token")
				}

				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if !ok {
					t.Error("failed to parse claims")
				}
				if claims["userId"] != "123" {
					t.Errorf("expected userId to be 123, got %v", claims["userId"])
				}
				if claims["expiredAt"] == nil {
					t.Error("expected expiredAt to be set")
				}
			},
		},
		{
			name:    "empty secret",
			secret:  []byte(""),
			userId:  123,
			wantErr: false,
			validate: func(t *testing.T, token string) {
				if token == "" {
					t.Error("expected non-empty token")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set a fixed expiration time for testing
			config.Envs.JWTExpiration = 3600 // 1 hour

			token, err := CreateJWT(tt.secret, tt.userId)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			tt.validate(t, token)
		})
	}
}
