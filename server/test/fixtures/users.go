package fixtures

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/example/go-webdb-template/internal/model"
	"github.com/example/go-webdb-template/internal/service"
)

// CreateTestUser creates a test user using the service layer
func CreateTestUser(t *testing.T, svc *service.UserService, name string) *model.User {
	req := &model.CreateUserRequest{
		Name:  name,
		Email: name + "@example.com",
	}
	user, err := svc.CreateUser(context.Background(), req)
	require.NoError(t, err)
	return user
}

// CreateTestUserWithEmail creates a test user with a specific email
func CreateTestUserWithEmail(t *testing.T, svc *service.UserService, name, email string) *model.User {
	req := &model.CreateUserRequest{
		Name:  name,
		Email: email,
	}
	user, err := svc.CreateUser(context.Background(), req)
	require.NoError(t, err)
	return user
}

// CreateMultipleTestUsers creates multiple test users
func CreateMultipleTestUsers(t *testing.T, svc *service.UserService, count int) []*model.User {
	users := make([]*model.User, count)
	for i := 0; i < count; i++ {
		name := "User" + string(rune('A'+i))
		users[i] = CreateTestUser(t, svc, name)
	}
	return users
}
