package utils

import (
	"net/http"
	"strings"
	"time"

	"Alice-Seahat-Healthcare/seahat-be/config"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type JWTClaims struct {
	Data     any
	Duration time.Duration
}

func JwtGenerate(claims JWTClaims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": claims.Data,
		"iss":  constant.JwtIssuer,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(claims.Duration).Unix(),
	})

	signed, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	return signed, nil
}

func JwtParse(signed string, secretKey string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(signed, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	},
		jwt.WithIssuer(constant.JwtIssuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, true
	}

	return nil, false
}

func JwtGenerateDefault(data any, secretKey string) (string, error) {
	claims := JWTClaims{
		Data:     data,
		Duration: constant.JwtDefaultDuration,
	}

	return JwtGenerate(claims, secretKey)
}

func JwtGenerateUser(data any) (string, error) {
	return JwtGenerateDefault(data, config.Jwt.SecretKey)
}

func JwtGenerateDoctor(data any) (string, error) {
	return JwtGenerateDefault(data, config.Jwt.DoctorKey)
}

func JwtGenerateManager(data any) (string, error) {
	return JwtGenerateDefault(data, config.Jwt.ManagerKey)
}

func JwtGenerateAdmin(data any) (string, error) {
	return JwtGenerateDefault(data, config.Jwt.AdminKey)
}

func JwtParseUser(signed string) (jwt.MapClaims, bool) {
	return JwtParse(signed, config.Jwt.SecretKey)
}

func JwtParseDoctor(signed string) (jwt.MapClaims, bool) {
	return JwtParse(signed, config.Jwt.DoctorKey)
}

func JwtParseManager(signed string) (jwt.MapClaims, bool) {
	return JwtParse(signed, config.Jwt.ManagerKey)
}

func JwtParseAdmin(signed string) (jwt.MapClaims, bool) {
	return JwtParse(signed, config.Jwt.AdminKey)
}

func GetUserFromJwt(req *http.Request) *entity.User {
	id, email := getActorFromJwt(req, JwtParseUser)

	if id == 0 || email == "" {
		return nil
	}

	return &entity.User{
		ID:    id,
		Email: email,
	}
}

func getActorFromJwt(req *http.Request, jwtFunc func(string) (jwt.MapClaims, bool)) (uint, string) {
	authorization := req.Header.Get("Authorization")
	bearerToken := strings.Split(authorization, " ")

	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		return 0, ""
	}

	user, ok := jwtFunc(bearerToken[1])
	if !ok {
		return 0, ""
	}

	userData, ok := user["data"].(map[string]any)
	if !ok {
		return 0, ""
	}

	id, ok := userData["ID"].(float64)
	if !ok {
		return 0, ""
	}

	email, ok := userData["Email"].(string)
	if !ok {
		return 0, ""
	}

	return uint(id), email
}
