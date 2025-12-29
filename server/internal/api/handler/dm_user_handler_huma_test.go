package handler

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

// TestRegisterUserEndpointsExists はRegisterUserEndpoints関数が存在することを確認
func TestRegisterUserEndpointsExists(t *testing.T) {
	// RegisterUserEndpoints関数のシグネチャを確認
	var _ func(api huma.API, h *DmUserHandler) = RegisterUserEndpoints
}
