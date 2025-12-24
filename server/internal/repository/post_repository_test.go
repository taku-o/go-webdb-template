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

func setupPostTestDB(t *testing.T) *sql.DB {
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

		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = db.Exec(schema)
	require.NoError(t, err)

	// Insert a test user
	_, err = db.Exec("INSERT INTO users (id, name, email) VALUES (?, ?, ?)", 1, "Test User", "test@example.com")
	require.NoError(t, err)

	return db
}

func TestPostRepository_Create(t *testing.T) {
	db := setupPostTestDB(t)
	defer db.Close()

	repo := repository.NewPostRepository(nil, nil)

	post := &model.Post{
		UserID:  1,
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	err := repo.Create(db, post)
	assert.NoError(t, err)
	assert.NotZero(t, post.ID)
	assert.NotZero(t, post.CreatedAt)
	assert.NotZero(t, post.UpdatedAt)
}

func TestPostRepository_GetByID(t *testing.T) {
	db := setupPostTestDB(t)
	defer db.Close()

	repo := repository.NewPostRepository(nil, nil)

	// Insert test data
	result, err := db.Exec(
		"INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
		1, "Test Post", "Test content",
	)
	require.NoError(t, err)
	id, _ := result.LastInsertId()

	// Test retrieval
	post, err := repo.GetByID(db, id)
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, int64(id), post.ID)
	assert.Equal(t, int64(1), post.UserID)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "Test content", post.Content)
}

func TestPostRepository_GetByID_NotFound(t *testing.T) {
	db := setupPostTestDB(t)
	defer db.Close()

	repo := repository.NewPostRepository(nil, nil)

	// Test retrieval of non-existent post
	post, err := repo.GetByID(db, 999)
	assert.Error(t, err)
	assert.Nil(t, post)
}

func TestPostRepository_Update(t *testing.T) {
	db := setupPostTestDB(t)
	defer db.Close()

	repo := repository.NewPostRepository(nil, nil)

	// Insert test data
	result, err := db.Exec(
		"INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
		1, "Original Title", "Original content",
	)
	require.NoError(t, err)
	id, _ := result.LastInsertId()

	// Update post
	post := &model.Post{
		ID:      id,
		UserID:  1,
		Title:   "Updated Title",
		Content: "Updated content",
	}

	err = repo.Update(db, post)
	assert.NoError(t, err)

	// Verify update
	updated, err := repo.GetByID(db, id)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, "Updated content", updated.Content)
}

func TestPostRepository_Delete(t *testing.T) {
	db := setupPostTestDB(t)
	defer db.Close()

	repo := repository.NewPostRepository(nil, nil)

	// Insert test data
	result, err := db.Exec(
		"INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
		1, "Test Post", "Test content",
	)
	require.NoError(t, err)
	id, _ := result.LastInsertId()

	// Delete post
	err = repo.Delete(db, id)
	assert.NoError(t, err)

	// Verify deletion
	post, err := repo.GetByID(db, id)
	assert.Error(t, err)
	assert.Nil(t, post)
}

func TestPostRepository_GetByUserID(t *testing.T) {
	db := setupPostTestDB(t)
	defer db.Close()

	repo := repository.NewPostRepository(nil, nil)

	// Insert test posts
	_, err := db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", 1, "Post 1", "Content 1")
	require.NoError(t, err)
	_, err = db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", 1, "Post 2", "Content 2")
	require.NoError(t, err)

	// Get posts by user ID
	posts, err := repo.GetByUserID(db, 1)
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.Equal(t, "Post 1", posts[0].Title)
	assert.Equal(t, "Post 2", posts[1].Title)
}
