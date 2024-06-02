package utils

import (
	"Alice-Seahat-Healthcare/seahat-be/config"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string, cost int) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), cost)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return hash, nil
}

func HashPasswordDefault(pwd *string) (*string, error) {
	if pwd == nil {
		return nil, nil
	}

	hashedByte, err := HashPassword(*pwd, config.Hash.Cost)
	if err != nil {
		return nil, err
	}

	hashedPass := string(hashedByte)

	return &hashedPass, nil
}

func HashComparePassword(pwd string, hash []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hash, []byte(pwd)); err != nil {
		return false
	}

	return true
}

func HashCompareDefault(pwd string, hashPwd *string) bool {
	password := ""
	if hashPwd != nil {
		password = *hashPwd
	}

	return HashComparePassword(pwd, []byte(password))
}
