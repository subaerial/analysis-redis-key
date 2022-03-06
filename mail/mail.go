package mail

import (
	"analysis.redis/config"
	"gopkg.in/gomail.v2"
)

func buildMsg(from string, to []string, subject string, body string) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to...)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)
	return message
}

func sendEmail(from string, to []string, subject string, body string) {
	dialer := gomail.Dialer{
		Host:     config.Properties.Mail.Dialer.Host,
		Port:     config.Properties.Mail.Dialer.Port,
		Username: config.Properties.Mail.Dialer.Username,
		Password: config.Properties.Mail.Dialer.Password,
	}
	if err := dialer.DialAndSend(buildMsg(from, to, subject, body)); err != nil {
		panic(err)
	}
}
