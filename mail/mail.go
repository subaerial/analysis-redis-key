package mail

import "gopkg.in/gomail.v2"

func buildMsg(from string, to string, subject string, body string) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)
	return message
}

func SendEmail(from string, to string, subject string, body string) {
	dialer := gomail.Dialer{
		Host:     "smtp.126.com",
		Port:     25,
		Username: "moxfan@126.com",
		Password: "stable110",
	}
	if err := dialer.DialAndSend(buildMsg(from, to, subject, body)); err != nil {
		panic(err)
	}
}
