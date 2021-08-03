package mailer

import (
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Service interface {
	SendEmail(toEmail string, toName string, subject string, plainBody string, htmlBody string, options ...SendOption) error
}

type sendGrid struct {
	APIKey             string
	DefaultSenderName  string
	DefaultSenderEmail string
}

func NewSendGrid(APIKey, DefaultSenderName, DefaultSenderEmail string) Service {
	return &sendGrid{
		APIKey:             APIKey,
		DefaultSenderName:  DefaultSenderName,
		DefaultSenderEmail: DefaultSenderEmail,
	}
}

func (s sendGrid) SendEmail(toEmail string, toName string, subject string, plaintBody string, htmlBody string, options ...SendOption) error {
	cfg := &SendConfig{
		SenderName:     s.DefaultSenderName,
		SenderEmail:    s.DefaultSenderEmail,
		SenderTemplate: "",
	}
	for _, option := range options {
		option(cfg)
	}
	from := mail.NewEmail(cfg.SenderName, cfg.SenderEmail)
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, plaintBody, htmlBody)
	client := sendgrid.NewSendClient(s.APIKey)
	_, err := client.Send(message)
	return errors.Wrap(err, "failed to send verification email")
}
