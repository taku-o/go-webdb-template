package handler

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

// TestRegisterDmPostEndpointsExists はRegisterDmPostEndpoints関数が存在することを確認
func TestRegisterDmPostEndpointsExists(t *testing.T) {
	// RegisterDmPostEndpoints関数のシグネチャを確認
	var _ func(api huma.API, h *DmPostHandler) = RegisterDmPostEndpoints
}
