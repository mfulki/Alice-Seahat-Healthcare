package utils

import (
	"bytes"
	"os"
	"text/template"

	"Alice-Seahat-Healthcare/seahat-be/config"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/libs/mail"

	"github.com/sirupsen/logrus"
)

type htmlTemplate struct {
	fileName string
	data     map[string]string
}

func templateExecute(ht htmlTemplate) (string, error) {
	wd, _ := os.Getwd()
	htmlPath := wd + "/assets/html/mail/" + ht.fileName

	tmpl, err := template.ParseFiles(htmlPath)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, ht.data); err != nil {
		logrus.Error(err)
		return "", err
	}

	return buf.String(), nil
}

func sendEmailVerification(dm mail.MailDialer, email, token, url string) error {
	content, err := templateExecute(htmlTemplate{
		fileName: "verification.html",
		data: map[string]string{
			"verificationURL": url + "?token=" + token,
		},
	})

	if err != nil {
		return err
	}

	mm := mail.MailMessage{
		To:          []string{email},
		Subject:     constant.MailSubjectVerification,
		ContentHTML: content,
	}

	if err := dm.SendMessage(mm); err != nil {
		return err
	}

	return nil
}

func SendEmailVerification(dm mail.MailDialer, email, token string) error {
	return sendEmailVerification(dm, email, token, config.FE.VerificationURL)
}

func SendEmailVerificationDoctor(dm mail.MailDialer, email, token string) error {
	return sendEmailVerification(dm, email, token, config.FE.DoctorVerificationURL)
}

func SendEmailForgotToken(dm mail.MailDialer, email string, token string) error {
	content, err := templateExecute(htmlTemplate{
		fileName: "reset.html",
		data: map[string]string{
			"resetURL": config.FE.ResetURL + "?token=" + token,
		},
	})

	if err != nil {
		return err
	}

	mm := mail.MailMessage{
		To:          []string{email},
		Subject:     constant.MailSubjectReset,
		ContentHTML: content,
	}

	if err := dm.SendMessage(mm); err != nil {
		return err
	}

	return nil
}

func SendEmailAddPartner(dm mail.MailDialer, p entity.Partner, pwd string) error {
	content, err := templateExecute(htmlTemplate{
		fileName: "addPartner.html",
		data: map[string]string{
			"name":            p.Name,
			"managerName":     p.PharmacyManager.Name,
			"managerEmail":    p.PharmacyManager.Email,
			"managerPassword": pwd,
			"loginURL":        config.FE.ManagerLoginURL,
		},
	})

	if err != nil {
		return err
	}

	mm := mail.MailMessage{
		To:          []string{p.PharmacyManager.Email},
		Subject:     constant.MailSubjectNewPartner,
		ContentHTML: content,
	}

	if err := dm.SendMessage(mm); err != nil {
		return err
	}

	return nil
}
