package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var App = new(AppEnv)
var Hash = new(HashEnv)
var Jwt = new(JwtEnv)
var DB = new(DatabaseEnv)
var FE = new(FeEnv)
var SMTP = new(SmtpEnv)
var Cloudinary = new(UploadCloudinaryEnv)
var RajaOngkir = new(RajaOngkirEnv)

func Load() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatal(err)
	}

	if err := App.loadEnv(); err != nil {
		logrus.Fatal(err)
	}

	if err := FE.loadEnv(); err != nil {
		logrus.Fatal(err)
	}

	if err := Hash.loadEnv(); err != nil {
		logrus.Fatal(err)
	}

	if err := Jwt.loadEnv(); err != nil {
		logrus.Fatal(err)
	}

	if err := DB.loadEnv(); err != nil {
		logrus.Fatal(err)
	}

	if err := SMTP.loadEnv(); err != nil {
		logrus.Fatal(err)
	}

	if err := Cloudinary.loadEnv(); err != nil {
		logrus.Fatal(err)
	}

	if err := RajaOngkir.loadEnv(); err != nil {
		logrus.Fatal(err)
	}
}

func getEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("%s env is must be filled out", key)
	}

	return val, nil
}

func getIntEnv(key string) (int, error) {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0, fmt.Errorf("%s env is must be integer", key)
	}

	return val, nil
}
