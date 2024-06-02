package config

type SmtpEnv struct {
	Host     string
	Port     int
	Username string
	Password string
	MailFrom string
}

func (e *SmtpEnv) loadEnv() error {
	host, err := getEnv("SMTP_HOST")
	if err != nil {
		return err
	}

	port, err := getIntEnv("SMTP_PORT")
	if err != nil {
		return err
	}

	username, err := getEnv("SMTP_USERNAME")
	if err != nil {
		return err
	}

	password, err := getEnv("SMTP_PASSWORD")
	if err != nil {
		return err
	}

	mailForm, err := getEnv("SMTP_MAIL_FROM")
	if err != nil {
		return err
	}

	e.Host = host
	e.Port = port
	e.Username = username
	e.Password = password
	e.MailFrom = mailForm

	return nil
}
