package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/danielgtaylor/huma/v2"
	humaapi "github.com/taku-o/go-webdb-template/internal/api/huma"
	"github.com/taku-o/go-webdb-template/internal/auth"
	"github.com/taku-o/go-webdb-template/internal/usecase"
)

// EmailHandler はメール送信APIのハンドラー
type EmailHandler struct {
	emailUsecase *usecase.EmailUsecase
}

// NewEmailHandler は新しいEmailHandlerを作成
func NewEmailHandler(emailUsecase *usecase.EmailUsecase) *EmailHandler {
	return &EmailHandler{
		emailUsecase: emailUsecase,
	}
}

// RegisterEmailEndpoints はメール送信エンドポイントを登録
func RegisterEmailEndpoints(api huma.API, h *EmailHandler) {
	// POST /api/email/send - メール送信
	huma.Register(api, huma.Operation{
		OperationID:   "send-email",
		Method:        http.MethodPost,
		Path:          "/api/email/send",
		Summary:       "メールを送信",
		Description:   "**Access Level:** `public` (Public API Key JWT または Auth0 JWT でアクセス可能)",
		Tags:          []string{"email"},
		DefaultStatus: http.StatusOK,
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *humaapi.SendEmailInput) (*humaapi.SendEmailOutput, error) {
		// 公開レベルのチェック（publicエンドポイント）
		if err := auth.CheckAccessLevel(ctx, auth.AccessLevelPublic); err != nil {
			return nil, huma.Error403Forbidden(err.Error())
		}

		// メールアドレスのバリデーション
		for _, addr := range input.Body.To {
			if err := validateEmail(addr); err != nil {
				return nil, huma.Error400BadRequest(fmt.Sprintf("invalid email address: %s", addr))
			}
		}

		// usecase層でメール送信を実行
		if err := h.emailUsecase.SendEmail(ctx, input.Body.To, input.Body.Template, input.Body.Data); err != nil {
			return nil, huma.Error400BadRequest(fmt.Sprintf("failed to send email: %v", err))
		}

		resp := &humaapi.SendEmailOutput{}
		resp.Body.Success = true
		resp.Body.Message = "メールを送信しました"
		return resp, nil
	})
}

// validateEmail はメールアドレスの形式を検証
func validateEmail(addr string) error {
	if addr == "" {
		return fmt.Errorf("email address is empty")
	}
	_, err := mail.ParseAddress(addr)
	return err
}
