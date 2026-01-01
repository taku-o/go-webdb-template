package email

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// SESSender はAWS SESにメールを送信する実装
type SESSender struct {
	client *ses.Client
	from   string
}

// NewSESSender は新しいSESSenderを作成
func NewSESSender(region, from string) (*SESSender, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	client := ses.NewFromConfig(cfg)

	return &SESSender{
		client: client,
		from:   from,
	}, nil
}

// Send はAWS SES SDKを使用してメールを送信
func (s *SESSender) Send(ctx context.Context, to []string, subject, body string) error {
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: to,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(body),
				},
			},
		},
		Source: aws.String(s.from),
	}

	_, err := s.client.SendEmail(ctx, input)
	return err
}
