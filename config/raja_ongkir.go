package config

type RajaOngkirEnv struct {
	URL string
	Key string
}

func (e *RajaOngkirEnv) loadEnv() error {
	url, err := getEnv("RAJAONGKIR_URL")
	if err != nil {
		return err
	}

	key, err := getEnv("RAJAONGKIR_KEY")
	if err != nil {
		return err
	}

	e.URL = url
	e.Key = key

	return nil
}
