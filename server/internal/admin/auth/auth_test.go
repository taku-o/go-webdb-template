package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty string")
	}

	if hash == password {
		t.Error("Hash should not be the same as password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// 正しいパスワードでチェック
	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash should return true for correct password")
	}

	// 間違ったパスワードでチェック
	if CheckPasswordHash("wrongpassword", hash) {
		t.Error("CheckPasswordHash should return false for wrong password")
	}
}

func TestGetSessionLifetime(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"positive value", 3600, 3600},
		{"zero value", 0, 7200},
		{"negative value", -1, 7200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSessionLifetime(tt.input)
			if result != tt.expected {
				t.Errorf("GetSessionLifetime(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}
