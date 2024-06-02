package config

type DatabaseEnv struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func (e *DatabaseEnv) loadEnv() error {
	host, err := getEnv("DB_HOST")
	if err != nil {
		return err
	}

	port, err := getIntEnv("DB_PORT")
	if err != nil {
		return err
	}

	user, err := getEnv("DB_USER")
	if err != nil {
		return err
	}

	password, err := getEnv("DB_PASSWORD")
	if err != nil {
		return err
	}

	name, err := getEnv("DB_NAME")
	if err != nil {
		return err
	}

	e.Host = host
	e.Port = port
	e.User = user
	e.Password = password
	e.Name = name

	return nil
}
