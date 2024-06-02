package config

type AppEnv struct {
	Env  string
	Port int
}

func (e *AppEnv) loadEnv() error {
	enviroment, err := getEnv("APP_ENV")
	if err != nil {
		return err
	}

	port, err := getIntEnv("APP_PORT")
	if err != nil {
		return err
	}

	e.Env = enviroment
	e.Port = port

	return nil
}
