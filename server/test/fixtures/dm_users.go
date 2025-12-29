package fixtures

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/service"
)

// CreateTestUser creates a test user using the service layer
func CreateTestUser(t *testing.T, svc *service.DmUserService, name string) *model.DmUser {
	req := &model.CreateDmUserRequest{
		Name:  name,
		Email: name + "@example.com",
	}
	user, err := svc.CreateDmUser(context.Background(), req)
	require.NoError(t, err)
	return user
}

// CreateTestUserWithEmail creates a test user with a specific email
func CreateTestUserWithEmail(t *testing.T, svc *service.DmUserService, name, email string) *model.DmUser {
	req := &model.CreateDmUserRequest{
		Name:  name,
		Email: email,
	}
	user, err := svc.CreateDmUser(context.Background(), req)
	require.NoError(t, err)
	return user
}

// CreateMultipleTestUsers creates multiple test users
func CreateMultipleTestUsers(t *testing.T, svc *service.DmUserService, count int) []*model.DmUser {
	users := make([]*model.DmUser, count)
	for i := 0; i < count; i++ {
		name := "User" + string(rune('A'+i))
		users[i] = CreateTestUser(t, svc, name)
	}
	return users
}
