package internal

import (
	"notification-service/infra/config"
	"strconv"

	"github.com/go-mail/mail/v2"
)

type MailerImplementation interface {
	NewMessage() *mail.Message
	Dialer(envs *config.Env) *mail.Dialer
	MailPrepared(msg *mail.Message, envs *config.Env)
}

type MailerContent struct {
	To struct {
		Email string
	}
	Subject string
	Body    string
}

type Mailer struct {
	content MailerContent
}

func NewMailer(content MailerContent) MailerImplementation {
	return &Mailer{
		content: content,
	}
}

func (m *Mailer) NewMessage() *mail.Message {
	return mail.NewMessage()
}

func (m *Mailer) Dialer(envs *config.Env) *mail.Dialer {
	enumeredPort, _ := strconv.Atoi(envs.SmtpConfig.Port)

	return mail.NewDialer(
		envs.SmtpConfig.Host,
		enumeredPort,
		envs.SmtpConfig.EmailNoReplyUser,
		envs.SmtpConfig.EmailNoReplyPassword,
	)
}

func (m *Mailer) MailPrepared(msg *mail.Message, envs *config.Env) {
	msg.SetHeader("From", envs.SmtpConfig.EmailNoReply)
	msg.SetHeader("To", m.content.To.Email)
	msg.SetHeader("Subject", m.content.Subject)
	msg.SetBody("text/html", m.content.Body)
}
