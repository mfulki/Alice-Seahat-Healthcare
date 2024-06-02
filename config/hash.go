package config

type HashEnv struct {
	Cost int
}

func (e *HashEnv) loadEnv() error {
	hashCost, err := getIntEnv("HASH_COST")
	if err != nil {
		return err
	}

	e.Cost = hashCost

	return nil
}
