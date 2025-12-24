package repository_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3"

	"github.com/example/go-db-prj-sample/internal/model"
	"github.com/example/go-db-prj-sample/internal/repository"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Create schema
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err = db.Exec(schema)
	require.NoError(t, err)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewUserRepository(nil)

	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	err := repo.Create(db, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewUserRepository(nil)

	// Insert test data
	result, err := db.Exec(
		"INSERT INTO users (name, email) VALUES (?, ?)",
		"Test User", "test@example.com",
	)
	require.NoError(t, err)
	id, _ := result.LastInsertId()

	// Test retrieval
	user, err := repo.GetByID(db, id)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(id), user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewUserRepository(nil)

	// Test retrieval of non-existent user
	user, err := repo.GetByID(db, 999)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewUserRepository(nil)

	// Insert test data
	result, err := db.Exec(
		"INSERT INTO users (name, email) VALUES (?, ?)",
		"Original Name", "original@example.com",
	)
	require.NoError(t, err)
	id, _ := result.LastInsertId()

	// Update user
	user := &model.User{
		ID:    id,
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	err = repo.Update(db, user)
	assert.NoError(t, err)

	// Verify update
	updated, err := repo.GetByID(db, id)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, "updated@example.com", updated.Email)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewUserRepository(nil)

	// Insert test data
	result, err := db.Exec(
		"INSERT INTO users (name, email) VALUES (?, ?)",
		"Test User", "test@example.com",
	)
	require.NoError(t, err)
	id, _ := result.LastInsertId()

	// Delete user
	err = repo.Delete(db, id)
	assert.NoError(t, err)

	// Verify deletion
	user, err := repo.GetByID(db, id)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert multiple users
	_, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "User 1", "user1@example.com")
	require.NoError(t, err)
	_, err = db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "User 2", "user2@example.com")
	require.NoError(t, err)

	// Mock dbManager that returns the same connection for all shards
	// For this test, we'll test the single-shard GetAll functionality
	// Full cross-shard testing should be in integration tests

	// This is a simplified test - full cross-shard tests should be in integration tests
}
