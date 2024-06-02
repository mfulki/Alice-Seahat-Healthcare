package utils

import (
	"crypto/rand"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) (string, error) {
	random := make([]byte, length)
	_, err := rand.Read(random)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	charsetLength := byte(len(charset))
	for index, r := range random {
		random[index] = charset[r%charsetLength]
	}

	return string(random), nil
}

func DomainURL(uri string) string {
	url, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}

	hostname := strings.TrimPrefix(url.Hostname(), "www.")
	return hostname
}

func Geo2LongLat(geo string) (float64, float64, error) {
	parts := strings.Split(geo, "(")
	coords := strings.TrimRight(parts[1], ")")
	splitted := strings.Split(coords, " ")

	longitude, err := strconv.ParseFloat(splitted[0], 64)
	if err != nil {
		return 0, 0, err
	}

	latitude, err := strconv.ParseFloat(splitted[1], 64)
	if err != nil {
		return 0, 0, err
	}

	return longitude, latitude, nil
}
