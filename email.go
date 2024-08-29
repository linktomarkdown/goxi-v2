package goxi_v2

import "gopkg.in/gomail.v2"

type EmailLogic struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailLogic(params *EmailLogic) *EmailLogic {
	return &EmailLogic{
		Host:     params.Host,
		Port:     params.Port,
		Username: params.Username,
		Password: params.Password,
	}
}

// Send 发送邮件
func (e *EmailLogic) Send(to string, subject string, body string, attach string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	m.Attach(attach)

	d := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)
	if err := d.DialAndSend(m); err != nil {
		println(err)
		return nil
	}
	return nil
}
