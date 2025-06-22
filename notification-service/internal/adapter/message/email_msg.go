package message

import (
	"crypto/tls"

	"github.com/go-mail/mail"
	"github.com/labstack/gommon/log"

	"notification-service/config"
)

type MessageEmailInterfface interface {
	SendEmailNotif(form, subject, body string) error
}

type emailAttribute struct {
	Username string
	Password string
	Host     string
	Port     int
	From     string
	IsTLS    bool
}

// SendEmailNotif implements MessageEmailInterfface.
func (e *emailAttribute) SendEmailNotif(to string, subject string, body string) error {

	m := mail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to)

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := mail.Dialer{
		Host:     e.Host,
		Port:     e.Port,
		Username: e.Username,
		Password: e.Password,
	}

	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: e.IsTLS,
	}

	if err := d.DialAndSend(m); err != nil {
		log.Errorf("[SendEmailNotif-1] error: %v", err)
		return err
	}

	return nil
}

func NewMessageEmail(cfg *config.Config) MessageEmailInterfface {
	return &emailAttribute{
		Username: cfg.Email.Username,
		Password: cfg.Email.Password,
		Host:     cfg.Email.Host,
		Port:     cfg.Email.Port,
		From:     cfg.Email.Sending,
		IsTLS:    cfg.Email.IsTLS,
	}
}
