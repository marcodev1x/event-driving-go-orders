package internal

import (
	"bytes"
	"notification-service/infra/config"
	"strconv"
	"text/template"

	"github.com/go-mail/mail/v2"
)

type MailerImplementation interface {
	NewMessage() *mail.Message
	Dialer(envs *config.Env) *mail.Dialer
	MailPrepared(msg *mail.Message, envs *config.Env)
	SetupHtml(params interface{}, file string) error
}

type MailerContent struct {
	To      string
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
	msg.SetHeader("To", m.content.To)
	msg.SetHeader("Subject", m.content.Subject)
	msg.SetBody("text/html", m.content.Body)
}

func (m *Mailer) SetupHtml(params interface{}, file string) error {
	tpl, err := template.ParseFiles(file)

	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, params); err != nil {
		return err
	}

	m.content.Body = buf.String()

	return nil
}
