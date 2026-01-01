package email

import (
	"bytes"
	"fmt"
	"text/template"
)

// TemplateService はメールテンプレートサービス
type TemplateService struct {
	templates map[string]*template.Template
	subjects  map[string]string
}

// NewTemplateService は新しいTemplateServiceを作成
// テンプレートをソースコードに直書きで定義します
func NewTemplateService() *TemplateService {
	ts := &TemplateService{
		templates: make(map[string]*template.Template),
		subjects:  make(map[string]string),
	}

	// welcomeテンプレートの定義
	welcomeBody := `{{.Name}}様

go-webdb-templateへようこそ！

ご登録いただいたメールアドレス: {{.Email}}

ご不明な点がございましたら、お気軽にお問い合わせください。

よろしくお願いいたします。`

	welcomeTmpl, err := template.New("welcome").Parse(welcomeBody)
	if err != nil {
		panic(fmt.Sprintf("failed to parse welcome template: %v", err))
	}
	ts.templates["welcome"] = welcomeTmpl
	ts.subjects["welcome"] = "ようこそ go-webdb-template へ"

	return ts
}

// Render はテンプレート名とデータからメール本文を生成
func (s *TemplateService) Render(templateName string, data interface{}) (string, error) {
	tmpl, ok := s.templates[templateName]
	if !ok {
		return "", fmt.Errorf("template not found: %s", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// GetSubject はテンプレート名に基づいて件名を取得
func (s *TemplateService) GetSubject(templateName string) (string, error) {
	subject, ok := s.subjects[templateName]
	if !ok {
		return "", fmt.Errorf("subject not found for template: %s", templateName)
	}
	return subject, nil
}
