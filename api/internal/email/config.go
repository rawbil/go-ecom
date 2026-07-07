package email

import (
	"github.com/resend/resend-go/v3"
)

type ResendSender struct {
	client *resend.Client
	from   string
}

func NewResendSender(apikey, from string) *ResendSender {
	return &ResendSender{
		client: resend.NewClient(apikey),
		from:   from,
	}
}

func (s *ResendSender) EmailConfig(to, subject, html string) error {
	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: subject,
		Html:    html,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}
