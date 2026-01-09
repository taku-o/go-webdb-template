package fixtures

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/service"
	"github.com/taku-o/go-webdb-template/internal/util/idgen"
)

// CreateTestDmUser creates a test dm_user using the service layer
func CreateTestDmUser(t *testing.T, svc *service.DmUserService, name string) *model.DmUser {
	// ユニークなメールアドレスを生成
	uniqueID, err := idgen.GenerateUUIDv7()
	require.NoError(t, err)
	uniqueEmail := fmt.Sprintf("%s-%s@example.com", name, uniqueID)

	req := &model.CreateDmUserRequest{
		Name:  name,
		Email: uniqueEmail,
	}
	dmUser, err := svc.CreateDmUser(context.Background(), req)
	require.NoError(t, err)
	return dmUser
}

// CreateTestDmUserWithEmail creates a test dm_user with a specific email
func CreateTestDmUserWithEmail(t *testing.T, svc *service.DmUserService, name, email string) *model.DmUser {
	req := &model.CreateDmUserRequest{
		Name:  name,
		Email: email,
	}
	dmUser, err := svc.CreateDmUser(context.Background(), req)
	require.NoError(t, err)
	return dmUser
}

// CreateMultipleTestDmUsers creates multiple test dm_users
func CreateMultipleTestDmUsers(t *testing.T, svc *service.DmUserService, count int) []*model.DmUser {
	dmUsers := make([]*model.DmUser, count)
	for i := 0; i < count; i++ {
		name := "User" + string(rune('A'+i))
		dmUsers[i] = CreateTestDmUser(t, svc, name)
	}
	return dmUsers
}
