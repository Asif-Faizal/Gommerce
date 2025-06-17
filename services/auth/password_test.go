package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "successful password hash",
			password: "testPassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "long password (within bcrypt limit)",
			password: "veryLongPassword123456789012345678901234567890",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := HashPassword(tt.password)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if hashedPassword == "" {
				t.Error("expected non-empty hashed password")
			}
			if hashedPassword == tt.password {
				t.Error("hashed password should not match original password")
			}
		})
	}
}

func TestComparePasswords(t *testing.T) {
	// First create a hashed password to use in tests
	password := "testPassword123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		plainPassword  string
		want           bool
	}{
		{
			name:           "matching passwords",
			hashedPassword: hashedPassword,
			plainPassword:  password,
			want:           true,
		},
		{
			name:           "non-matching passwords",
			hashedPassword: hashedPassword,
			plainPassword:  "wrongPassword",
			want:           false,
		},
		{
			name:           "empty passwords",
			hashedPassword: hashedPassword,
			plainPassword:  "",
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComparePasswords(tt.hashedPassword, tt.plainPassword)
			if result != tt.want {
				t.Errorf("ComparePasswords() = %v, want %v", result, tt.want)
			}
		})
	}
}
