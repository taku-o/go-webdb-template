package api

import (
	"context"
)

// EmailServiceInterface はEmailServiceのインターフェース
type EmailServiceInterface interface {
	SendEmail(ctx context.Context, to []string, subject, body string) error
}

// TemplateServiceInterface はTemplateServiceのインターフェース
type TemplateServiceInterface interface {
	Render(templateName string, data any) (string, error)
	GetSubject(templateName string) (string, error)
}

// EmailUsecase はメール送信のビジネスロジックを担当するユースケース層
type EmailUsecase struct {
	emailService    EmailServiceInterface
	templateService TemplateServiceInterface
}

// NewEmailUsecase は新しいEmailUsecaseを作成
func NewEmailUsecase(emailService EmailServiceInterface, templateService TemplateServiceInterface) *EmailUsecase {
	return &EmailUsecase{
		emailService:    emailService,
		templateService: templateService,
	}
}

// SendEmail はテンプレートを使用してメールを送信
func (u *EmailUsecase) SendEmail(ctx context.Context, to []string, template string, data any) error {
	// テンプレートからメール本文を生成
	body, err := u.templateService.Render(template, data)
	if err != nil {
		return err
	}

	// テンプレートから件名を取得
	subject, err := u.templateService.GetSubject(template)
	if err != nil {
		return err
	}

	// メール送信
	return u.emailService.SendEmail(ctx, to, subject, body)
}
