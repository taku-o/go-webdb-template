package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateScope(t *testing.T) {
	tests := []struct {
		name    string
		scope   []string
		method  string
		wantErr bool
	}{
		{"read scope GET", []string{"read"}, "GET", false},
		{"read scope POST", []string{"read"}, "POST", true},
		{"write scope GET", []string{"write"}, "GET", true},
		{"write scope POST", []string{"write"}, "POST", false},
		{"write scope PUT", []string{"write"}, "PUT", false},
		{"write scope DELETE", []string{"write"}, "DELETE", false},
		{"read and write scope GET", []string{"read", "write"}, "GET", false},
		{"read and write scope POST", []string{"read", "write"}, "POST", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := &JWTClaims{Scope: tt.scope}
			err := validateScope(claims, tt.method)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
