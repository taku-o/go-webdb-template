package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestPrintUsersTSV(t *testing.T) {
	tests := []struct {
		name     string
		users    []*model.User
		wantRows int
		wantCols []string
	}{
		{
			name:     "empty users",
			users:    []*model.User{},
			wantRows: 1, // header only
			wantCols: []string{"ID", "Name", "Email", "CreatedAt", "UpdatedAt"},
		},
		{
			name: "single user",
			users: []*model.User{
				{
					ID:        1234567890123456789,
					Name:      "John Doe",
					Email:     "john@example.com",
					CreatedAt: time.Date(2025, 1, 27, 10, 30, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 27, 10, 30, 0, 0, time.UTC),
				},
			},
			wantRows: 2, // header + 1 user
			wantCols: []string{"ID", "Name", "Email", "CreatedAt", "UpdatedAt"},
		},
		{
			name: "multiple users",
			users: []*model.User{
				{
					ID:        1234567890123456789,
					Name:      "John Doe",
					Email:     "john@example.com",
					CreatedAt: time.Date(2025, 1, 27, 10, 30, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 27, 10, 30, 0, 0, time.UTC),
				},
				{
					ID:        1234567890123456790,
					Name:      "Jane Smith",
					Email:     "jane@example.com",
					CreatedAt: time.Date(2025, 1, 27, 11, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 27, 11, 0, 0, 0, time.UTC),
				},
			},
			wantRows: 3, // header + 2 users
			wantCols: []string{"ID", "Name", "Email", "CreatedAt", "UpdatedAt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			printUsersTSV(tt.users)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			// Split output into lines
			lines := strings.Split(strings.TrimSpace(output), "\n")

			// Check number of rows
			assert.Equal(t, tt.wantRows, len(lines), "unexpected number of rows")

			// Check header columns
			headerCols := strings.Split(lines[0], "\t")
			assert.Equal(t, tt.wantCols, headerCols, "unexpected header columns")

			// Check data rows
			for i, user := range tt.users {
				dataCols := strings.Split(lines[i+1], "\t")
				assert.Equal(t, 5, len(dataCols), "unexpected number of columns in data row")

				// Check ID
				expectedID := fmt.Sprintf("%d", user.ID)
				assert.Equal(t, expectedID, dataCols[0], "ID should match")

				// Check Name
				assert.Equal(t, user.Name, dataCols[1], "Name should match")

				// Check Email
				assert.Equal(t, user.Email, dataCols[2], "Email should match")

				// Check date format (RFC3339)
				assert.Contains(t, dataCols[3], "2025-01-27T", "CreatedAt should be in RFC3339 format")
				assert.Contains(t, dataCols[4], "2025-01-27T", "UpdatedAt should be in RFC3339 format")
			}
		})
	}
}

func TestPrintUsersTSV_RFC3339Format(t *testing.T) {
	users := []*model.User{
		{
			ID:        1234567890123456789,
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Date(2025, 1, 27, 10, 30, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 27, 15, 45, 30, 0, time.UTC),
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printUsersTSV(users)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Equal(t, 2, len(lines), "should have header and one data row")

	dataCols := strings.Split(lines[1], "\t")
	assert.Equal(t, "2025-01-27T10:30:00Z", dataCols[3], "CreatedAt should be RFC3339")
	assert.Equal(t, "2025-01-27T15:45:30Z", dataCols[4], "UpdatedAt should be RFC3339")
}

func TestValidateLimit(t *testing.T) {
	tests := []struct {
		name        string
		limit       int
		wantLimit   int
		wantErr     bool
		wantWarning bool
	}{
		{
			name:        "valid limit",
			limit:       20,
			wantLimit:   20,
			wantErr:     false,
			wantWarning: false,
		},
		{
			name:        "minimum valid limit",
			limit:       1,
			wantLimit:   1,
			wantErr:     false,
			wantWarning: false,
		},
		{
			name:        "maximum valid limit",
			limit:       100,
			wantLimit:   100,
			wantErr:     false,
			wantWarning: false,
		},
		{
			name:        "limit below minimum",
			limit:       0,
			wantLimit:   0,
			wantErr:     true,
			wantWarning: false,
		},
		{
			name:        "negative limit",
			limit:       -1,
			wantLimit:   0,
			wantErr:     true,
			wantWarning: false,
		},
		{
			name:        "limit above maximum",
			limit:       200,
			wantLimit:   100,
			wantErr:     false,
			wantWarning: true,
		},
		{
			name:        "limit at 101",
			limit:       101,
			wantLimit:   100,
			wantErr:     false,
			wantWarning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotErr, gotWarning := validateLimit(tt.limit)

			if tt.wantErr {
				assert.Error(t, gotErr, "expected error")
			} else {
				assert.NoError(t, gotErr, "unexpected error")
				assert.Equal(t, tt.wantLimit, gotLimit, "unexpected limit value")
			}

			assert.Equal(t, tt.wantWarning, gotWarning, "unexpected warning value")
		})
	}
}
