package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
	"github.com/taku-o/go-webdb-template/internal/util/idgen"
	"github.com/taku-o/go-webdb-template/test/testutil"
)

func TestDmUserRepository_Create(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueEmail := fmt.Sprintf("test-%s@example.com", uniqueID)

	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: uniqueEmail,
	}

	dmUser, err := dmUserRepo.Create(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, dmUser)

	// クリーンアップ
	defer func() {
		if dmUser != nil {
			_ = dmUserRepo.Delete(ctx, dmUser.ID)
		}
	}()

	assert.NotZero(t, dmUser.ID)
	assert.Equal(t, "Test User", dmUser.Name)
	assert.Equal(t, uniqueEmail, dmUser.Email)
	assert.NotZero(t, dmUser.CreatedAt)
	assert.NotZero(t, dmUser.UpdatedAt)
}

func TestDmUserRepository_GetByID(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueEmail := fmt.Sprintf("test-%s@example.com", uniqueID)

	// Create test user first
	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: uniqueEmail,
	}
	created, err := dmUserRepo.Create(ctx, req)
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmUserRepo.Delete(ctx, created.ID)
	}()

	// Test retrieval
	dmUser, err := dmUserRepo.GetByID(ctx, created.ID)
	assert.NoError(t, err)
	assert.NotNil(t, dmUser)
	assert.Equal(t, created.ID, dmUser.ID)
	assert.Equal(t, "Test User", dmUser.Name)
	assert.Equal(t, uniqueEmail, dmUser.Email)
}

func TestDmUserRepository_GetByID_NotFound(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// Test retrieval of non-existent user
	dmUser, err := dmUserRepo.GetByID(ctx, "00000000000000000000000000000000")
	assert.Error(t, err)
	assert.Nil(t, dmUser)
}

func TestDmUserRepository_Update(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	originalEmail := fmt.Sprintf("original-%s@example.com", uniqueID)
	updatedEmail := fmt.Sprintf("updated-%s@example.com", uniqueID)

	// Create test user first
	createReq := &model.CreateDmUserRequest{
		Name:  "Original Name",
		Email: originalEmail,
	}
	created, err := dmUserRepo.Create(ctx, createReq)
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmUserRepo.Delete(ctx, created.ID)
	}()

	// Update user
	updateReq := &model.UpdateDmUserRequest{
		Name:  "Updated Name",
		Email: updatedEmail,
	}
	updated, err := dmUserRepo.Update(ctx, created.ID, updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, updatedEmail, updated.Email)

	// Verify update
	dmUser, err := dmUserRepo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", dmUser.Name)
	assert.Equal(t, updatedEmail, dmUser.Email)
}

func TestDmUserRepository_Delete(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueEmail := fmt.Sprintf("test-%s@example.com", uniqueID)

	// Create test user first
	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: uniqueEmail,
	}
	created, err := dmUserRepo.Create(ctx, req)
	require.NoError(t, err)

	// Delete user
	err = dmUserRepo.Delete(ctx, created.ID)
	assert.NoError(t, err)

	// Verify deletion
	dmUser, err := dmUserRepo.GetByID(ctx, created.ID)
	assert.Error(t, err)
	assert.Nil(t, dmUser)
}

func TestDmUserRepository_CreateAndRetrieve(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID1, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueID2, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)

	// Create test users
	req1 := &model.CreateDmUserRequest{
		Name:  "User 1",
		Email: fmt.Sprintf("user1-%s@example.com", uniqueID1),
	}
	user1, err := dmUserRepo.Create(ctx, req1)
	require.NoError(t, err)

	req2 := &model.CreateDmUserRequest{
		Name:  "User 2",
		Email: fmt.Sprintf("user2-%s@example.com", uniqueID2),
	}
	user2, err := dmUserRepo.Create(ctx, req2)
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmUserRepo.Delete(ctx, user1.ID)
		_ = dmUserRepo.Delete(ctx, user2.ID)
	}()

	// Verify users can be retrieved by ID (single shard queries)
	retrieved1, err := dmUserRepo.GetByID(ctx, user1.ID)
	assert.NoError(t, err)
	assert.Equal(t, user1.ID, retrieved1.ID)
	assert.Equal(t, "User 1", retrieved1.Name)

	retrieved2, err := dmUserRepo.GetByID(ctx, user2.ID)
	assert.NoError(t, err)
	assert.Equal(t, user2.ID, retrieved2.ID)
	assert.Equal(t, "User 2", retrieved2.Name)
}

func TestDmUserRepository_InsertDmUsersBatch(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	tests := []struct {
		name      string
		tableName string
		dmUsers   []*model.DmUser
		wantErr   bool
	}{
		{
			name:      "empty slice",
			tableName: "dm_users_000",
			dmUsers:   []*model.DmUser{},
			wantErr:   false,
		},
		{
			name:      "single user",
			tableName: "dm_users_000",
			dmUsers: func() []*model.DmUser {
				id, _ := idgen.GenerateUUIDv7()
				return []*model.DmUser{
					{
						ID:    id,
						Name:  "Batch User 1",
						Email: fmt.Sprintf("batch1-%s@example.com", id),
					},
				}
			}(),
			wantErr: false,
		},
		{
			name:      "multiple users",
			tableName: "dm_users_000",
			dmUsers: func() []*model.DmUser {
				var users []*model.DmUser
				for i := 0; i < 3; i++ {
					id, _ := idgen.GenerateUUIDv7()
					users = append(users, &model.DmUser{
						ID:    id,
						Name:  fmt.Sprintf("Batch User %d", i+1),
						Email: fmt.Sprintf("batch%d-%s@example.com", i+1, id),
					})
				}
				return users
			}(),
			wantErr: false,
		},
		{
			name:      "invalid table name",
			tableName: "invalid_table",
			dmUsers: func() []*model.DmUser {
				id, _ := idgen.GenerateUUIDv7()
				return []*model.DmUser{
					{
						ID:    id,
						Name:  "Test User",
						Email: fmt.Sprintf("test-%s@example.com", id),
					},
				}
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dmUserRepo.InsertDmUsersBatch(ctx, tt.tableName, tt.dmUsers)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// クリーンアップ: 挿入されたデータを削除
				for _, u := range tt.dmUsers {
					_ = dmUserRepo.Delete(ctx, u.ID)
				}
			}
		})
	}
}


func TestDmUserRepository_CheckEmailExists(t *testing.T) {
	groupManager := testutil.SetupTestGroupManager(t, 4, 8)
	defer testutil.CleanupTestGroupManager(groupManager)

	dmUserRepo := repository.NewDmUserRepository(groupManager)
	ctx := context.Background()

	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	existingEmail := fmt.Sprintf("existing-%s@example.com", uniqueID)
	nonExistingEmail := fmt.Sprintf("non-existing-%s@example.com", uniqueID)

	// テスト用ユーザーを作成
	req := &model.CreateDmUserRequest{
		Name:  "Test User",
		Email: existingEmail,
	}
	created, err := dmUserRepo.Create(ctx, req)
	require.NoError(t, err)

	// クリーンアップ
	defer func() {
		_ = dmUserRepo.Delete(ctx, created.ID)
	}()

	tests := []struct {
		name    string
		email   string
		want    bool
		wantErr bool
	}{
		{
			name:    "メールアドレスが存在する場合",
			email:   existingEmail,
			want:    true,
			wantErr: false,
		},
		{
			name:    "メールアドレスが存在しない場合",
			email:   nonExistingEmail,
			want:    false,
			wantErr: false,
		},
		{
			name:    "空のメールアドレス",
			email:   "",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dmUserRepo.CheckEmailExists(ctx, tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
