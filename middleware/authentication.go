package middleware

import (
	"strings"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) Authentication(jwtFunc func(signed string) (jwt.MapClaims, bool)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		bearerToken := strings.Split(authorization, " ")

		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			ctx.Error(apperror.Unauthorized)
			ctx.Abort()
			return
		}

		user, ok := jwtFunc(bearerToken[1])
		if !ok {
			ctx.Error(apperror.Unauthorized)
			ctx.Abort()
			return
		}

		userDataMap, ok := user["data"].(map[string]any)
		if !ok {
			ctx.Error(apperror.Unauthorized)
			ctx.Abort()
			return
		}

		if role, ok := user["role"]; ok {
			userDataMap["role"] = role
		}

		ctx.Set(constant.UserContext, userDataMap)
	}
}

func (m *Middleware) AuthMulti(actors []string) gin.HandlerFunc {
	return m.Authentication(func(signed string) (jwt.MapClaims, bool) {
		for _, actor := range actors {
			if act, ok := m.isActorAuthed(actor, signed); ok {
				act["role"] = actor
				return act, ok
			}
		}

		return nil, false
	})
}

func (m *Middleware) UserAuth() gin.HandlerFunc {
	return m.Authentication(func(signed string) (jwt.MapClaims, bool) {
		user, ok := utils.JwtParseUser(signed)
		if !ok {
			return nil, false
		}

		user["role"] = constant.User
		return user, true
	})
}

func (m *Middleware) DoctorAuth() gin.HandlerFunc {
	return m.Authentication(func(signed string) (jwt.MapClaims, bool) {
		doctor, ok := utils.JwtParseDoctor(signed)
		if !ok {
			return nil, false
		}

		doctor["role"] = constant.Doctor
		return doctor, true
	})
}

func (m *Middleware) ManagerAuth() gin.HandlerFunc {
	return m.Authentication(func(signed string) (jwt.MapClaims, bool) {
		manager, ok := utils.JwtParseManager(signed)
		if !ok {
			return nil, false
		}

		manager["role"] = constant.Manager
		return manager, true
	})
}

func (m *Middleware) AdminAuth() gin.HandlerFunc {
	return m.Authentication(func(signed string) (jwt.MapClaims, bool) {
		admin, ok := utils.JwtParseAdmin(signed)
		if !ok {
			return nil, false
		}

		admin["role"] = constant.Admin
		return admin, true
	})
}

func (m *Middleware) isActorAuthed(actor string, signed string) (jwt.MapClaims, bool) {
	act, ok := make(jwt.MapClaims), false

	switch actor {
	case constant.User:
		act, ok = utils.JwtParseUser(signed)
	case constant.Doctor:
		act, ok = utils.JwtParseDoctor(signed)
	case constant.Manager:
		act, ok = utils.JwtParseManager(signed)
	case constant.Admin:
		act, ok = utils.JwtParseAdmin(signed)
	}

	return act, ok
}
