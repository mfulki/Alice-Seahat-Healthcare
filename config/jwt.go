package config

type JwtEnv struct {
	SecretKey  string
	DoctorKey  string
	ManagerKey string
	AdminKey   string
}

func (e *JwtEnv) loadEnv() error {
	secretKey, err := getEnv("JWT_SECRET_KEY")
	if err != nil {
		return err
	}

	doctorKey, err := getEnv("JWT_DOCTOR_SECRET_KEY")
	if err != nil {
		return err
	}

	managerKey, err := getEnv("JWT_PM_SECRET_KEY")
	if err != nil {
		return err
	}

	adminKey, err := getEnv("JWT_ADMIN_SECRET_KEY")
	if err != nil {
		return err
	}

	e.SecretKey = secretKey
	e.DoctorKey = doctorKey
	e.ManagerKey = managerKey
	e.AdminKey = adminKey

	return nil
}
