package service

import (
	"errors"
	"fmt"
	"net/smtp"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/graph/model"
)

type mailerService struct {
	Host           string
	Port           string
	User           string
	Pass           string
	RecipientEmail string
}

// NewMailerService creates and returns new mailer service
func NewMailerService(config string, recipientEmail string, validate *validator.Validate) (model.MailerService, error) {
	if err := validate.Var(recipientEmail, "required,email"); err != nil {
		return nil, errors.New("no valid email address for recipient")
	}

	mailerTransportRegex := regexp.MustCompile(`^smtp://(([a-zA-Z0-9+\-@]+):([^:@]+)@)?([a-zA-Z0-9\-.]*)(:([0-9]+))?$`)
	matches := mailerTransportRegex.FindStringSubmatch(config)
	if matches == nil {
		return nil, errors.New("could not parse configuration string")
	}

	ret := &mailerService{
		User:           matches[2],
		Pass:           matches[3],
		Host:           matches[4],
		Port:           matches[6],
		RecipientEmail: recipientEmail,
	}
	if ret.Port == "" {
		ret.Port = "25"
	}
	return ret, nil
}

// SendContactMail sends an email to the default contact email address
func (s *mailerService) SendContactMail(senderEmail string, senderName string, subject string, message string, ip string) error {
	var auth smtp.Auth
	if s.User != "" {
		auth = smtp.PlainAuth("", s.User, s.Pass, s.Host)
	}

	msg := []byte("To: " + s.RecipientEmail + "\r\n" +
		"From: \"" + senderName + "\" <" + senderEmail + ">\r\n" +
		"Subject: " + subject + "\r\n" +
		"X-Remote-IP: " + ip + "\r\n" +
		"\r\n" +
		message + "\r\n")

	err := smtp.SendMail(s.Host+":"+s.Port, auth, senderEmail, []string{s.RecipientEmail}, msg)

	if err != nil {
		return fmt.Errorf("could not send email (host: %s): %w", s.Host, err)
	}
	return nil
}
