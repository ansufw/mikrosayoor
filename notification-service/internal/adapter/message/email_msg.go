package message

import "notification-service/config"

type MessageEmailInterfface interface {
	SendEmailNotif(form, subject, body string)
}

type emailAttribute struct {
	Username string
	Password string
	Host     string
	Port     int
	From     string
}

// SendEmailNotif implements MessageEmailInterfface.
func (e *emailAttribute) SendEmailNotif(form string, subject string, body string) {
	panic("unimplemented")
}

func NewMessageEmail(cfg *config.Config) MessageEmailInterfface {
	return &emailAttribute{
		From:     cfg.Email.Host,
		Password: cfg.Email.Password,
		Host:     cfg.Email.Host,
		Port:     cfg.Email.Port,
	}
}
