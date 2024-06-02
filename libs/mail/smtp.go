package mail

import (
	"Alice-Seahat-Healthcare/seahat-be/config"
	"Alice-Seahat-Healthcare/seahat-be/constant"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type MailDialer interface {
	SendMessage(mm MailMessage) error
}

type dialer struct {
	mail *gomail.Dialer
	from string
}

func NewDialer() (*dialer, error) {
	d := &dialer{
		from: config.SMTP.MailFrom,
		mail: gomail.NewDialer(config.SMTP.Host, config.SMTP.Port, config.SMTP.Username, config.SMTP.Password),
	}

	if config.App.Env == constant.Production {
		if err := d.ping(); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (d *dialer) SendMessage(mm MailMessage) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", d.from)
	mail.SetHeader("To", mm.To...)
	mail.SetHeader("Subject", mm.Subject)
	mail.SetBody("text/html", mm.ContentHTML)

	if err := d.mail.DialAndSend(mail); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (d *dialer) ping() error {
	mail, err := d.mail.Dial()
	if err != nil {
		logrus.Error(err)
		return err
	}

	defer mail.Close()
	return nil
}
