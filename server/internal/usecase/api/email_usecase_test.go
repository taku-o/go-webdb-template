package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockEmailService はテスト用のEmailServiceモック
type MockEmailService struct {
	SendEmailFunc func(ctx context.Context, to []string, subject, body string) error
}

func (m *MockEmailService) SendEmail(ctx context.Context, to []string, subject, body string) error {
	if m.SendEmailFunc != nil {
		return m.SendEmailFunc(ctx, to, subject, body)
	}
	return nil
}

// MockTemplateService はテスト用のTemplateServiceモック
type MockTemplateService struct {
	RenderFunc     func(templateName string, data any) (string, error)
	GetSubjectFunc func(templateName string) (string, error)
}

func (m *MockTemplateService) Render(templateName string, data any) (string, error) {
	if m.RenderFunc != nil {
		return m.RenderFunc(templateName, data)
	}
	return "", nil
}

func (m *MockTemplateService) GetSubject(templateName string) (string, error) {
	if m.GetSubjectFunc != nil {
		return m.GetSubjectFunc(templateName)
	}
	return "", nil
}

func TestEmailUsecase_SendEmail(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name             string
		to               []string
		template         string
		data             any
		mockEmailService *MockEmailService
		mockTemplateService *MockTemplateService
		wantErr          bool
	}{
		{
			name:     "sends email successfully",
			to:       []string{"test@example.com"},
			template: "welcome",
			data:     map[string]any{"Name": "Test User"},
			mockEmailService: &MockEmailService{
				SendEmailFunc: func(ctx context.Context, to []string, subject, body string) error {
					return nil
				},
			},
			mockTemplateService: &MockTemplateService{
				RenderFunc: func(templateName string, data any) (string, error) {
					return "Hello Test User!", nil
				},
				GetSubjectFunc: func(templateName string) (string, error) {
					return "Welcome", nil
				},
			},
			wantErr: false,
		},
		{
			name:     "returns error when template render fails",
			to:       []string{"test@example.com"},
			template: "nonexistent",
			data:     map[string]any{"Name": "Test User"},
			mockEmailService: &MockEmailService{},
			mockTemplateService: &MockTemplateService{
				RenderFunc: func(templateName string, data any) (string, error) {
					return "", errors.New("template not found")
				},
			},
			wantErr: true,
		},
		{
			name:     "returns error when get subject fails",
			to:       []string{"test@example.com"},
			template: "welcome",
			data:     map[string]any{"Name": "Test User"},
			mockEmailService: &MockEmailService{},
			mockTemplateService: &MockTemplateService{
				RenderFunc: func(templateName string, data any) (string, error) {
					return "Hello Test User!", nil
				},
				GetSubjectFunc: func(templateName string) (string, error) {
					return "", errors.New("subject not found")
				},
			},
			wantErr: true,
		},
		{
			name:     "returns error when email send fails",
			to:       []string{"test@example.com"},
			template: "welcome",
			data:     map[string]any{"Name": "Test User"},
			mockEmailService: &MockEmailService{
				SendEmailFunc: func(ctx context.Context, to []string, subject, body string) error {
					return errors.New("send failed")
				},
			},
			mockTemplateService: &MockTemplateService{
				RenderFunc: func(templateName string, data any) (string, error) {
					return "Hello Test User!", nil
				},
				GetSubjectFunc: func(templateName string) (string, error) {
					return "Welcome", nil
				},
			},
			wantErr: true,
		},
		{
			name:     "sends email to multiple recipients",
			to:       []string{"user1@example.com", "user2@example.com", "user3@example.com"},
			template: "welcome",
			data:     map[string]any{"Name": "Users"},
			mockEmailService: &MockEmailService{
				SendEmailFunc: func(ctx context.Context, to []string, subject, body string) error {
					assert.Len(t, to, 3)
					return nil
				},
			},
			mockTemplateService: &MockTemplateService{
				RenderFunc: func(templateName string, data any) (string, error) {
					return "Hello Users!", nil
				},
				GetSubjectFunc: func(templateName string) (string, error) {
					return "Welcome", nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := NewEmailUsecase(tt.mockEmailService, tt.mockTemplateService)

			err := usecase.SendEmail(ctx, tt.to, tt.template, tt.data)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
