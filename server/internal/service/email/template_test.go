package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTemplateService(t *testing.T) {
	service := NewTemplateService()
	assert.NotNil(t, service)
}

func TestTemplateService_Render(t *testing.T) {
	service := NewTemplateService()

	tests := []struct {
		name         string
		templateName string
		data         interface{}
		wantErr      bool
		wantContains string
	}{
		{
			name:         "welcomeテンプレートが正しく置換される",
			templateName: "welcome",
			data: map[string]interface{}{
				"Name":  "田中太郎",
				"Email": "tanaka@example.com",
			},
			wantErr:      false,
			wantContains: "田中太郎",
		},
		{
			name:         "存在しないテンプレートでエラーを返す",
			templateName: "nonexistent",
			data:         map[string]interface{}{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Render(tt.templateName, tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, result, tt.wantContains)
			}
		})
	}
}

func TestTemplateService_GetSubject(t *testing.T) {
	service := NewTemplateService()

	tests := []struct {
		name         string
		templateName string
		wantErr      bool
		wantSubject  string
	}{
		{
			name:         "welcomeテンプレートの件名を取得",
			templateName: "welcome",
			wantErr:      false,
			wantSubject:  "ようこそ",
		},
		{
			name:         "存在しないテンプレートでエラーを返す",
			templateName: "nonexistent",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subject, err := service.GetSubject(tt.templateName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, subject, tt.wantSubject)
			}
		})
	}
}
