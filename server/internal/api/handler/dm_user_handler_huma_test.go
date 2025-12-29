package handler

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

// TestRegisterDmUserEndpointsExists はRegisterDmUserEndpoints関数が存在することを確認
func TestRegisterDmUserEndpointsExists(t *testing.T) {
	// RegisterDmUserEndpoints関数のシグネチャを確認
	var _ func(api huma.API, h *DmUserHandler) = RegisterDmUserEndpoints
}
