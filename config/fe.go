package config

type FeEnv struct {
	URL                   string
	ResetURL              string
	VerificationURL       string
	DoctorVerificationURL string
	ManagerLoginURL       string
}

func (e *FeEnv) loadEnv() error {
	url, err := getEnv("FE_URL")
	if err != nil {
		return err
	}

	verificationURL, err := getEnv("FE_VERIFICATION_URL")
	if err != nil {
		return err
	}

	verificationDoctorURL, err := getEnv("FE_DOCTOR_VERIFICATION_URL")
	if err != nil {
		return err
	}

	resetURL, err := getEnv("FE_RESET_URL")
	if err != nil {
		return err
	}

	pmLoginURL, err := getEnv("FE_PM_LOGIN_URL")
	if err != nil {
		return err
	}

	e.URL = url
	e.VerificationURL = verificationURL
	e.DoctorVerificationURL = verificationDoctorURL
	e.ResetURL = resetURL
	e.ManagerLoginURL = pmLoginURL

	return nil
}
