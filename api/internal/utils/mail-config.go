package utils

import (
	"errors"

	"github.com/rawbil/ecom2/internal/config"
	"github.com/resend/resend-go/v3"
)

type MailOptions struct {
	Html    string
	Subject string
	To      string
	To2 string
}

func SendMail(m MailOptions) error {
	api_key := config.GetResendConfig().ApiKey
	from := config.GetResendConfig().EmailFrom

	if api_key == "" {
		Log.Error("resend api key not found")
		return errors.New("resend api key not found")
	}
	if from == "" {
		Log.Error("resend email missing")
		return errors.New("resend email missing")
	}

	client := resend.NewClient(api_key)

	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{m.To, m.To2},
		Html:    m.Html,
		Subject: m.Subject,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		Log.Error("error sending email", "error", err)
		return err
	}

	Log.Info("Email sent successfully", "resend_id", sent.Id)

	return nil
}

// type ResendSender struct {
// 	client *resend.Client
// 	from   string
// }

// func NewResendSender(apikey, from string) *ResendSender {
// 	return &ResendSender{
// 		client: resend.NewClient(apikey),
// 		from:   from,
// 	}
// }

// func (s *ResendSender) EmailConfig(to, subject, html string) error {
// 	params := &resend.SendEmailRequest{
// 		From:    s.from,
// 		To:      []string{to},
// 		Subject: subject,
// 		Html:    html,
// 	}

// 	_, err := s.client.Emails.Send(params)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
