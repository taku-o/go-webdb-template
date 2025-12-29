package handler

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

// TestRegisterPostEndpointsExists はRegisterPostEndpoints関数が存在することを確認
func TestRegisterPostEndpointsExists(t *testing.T) {
	// RegisterPostEndpoints関数のシグネチャを確認
	var _ func(api huma.API, h *DmPostHandler) = RegisterPostEndpoints
}
