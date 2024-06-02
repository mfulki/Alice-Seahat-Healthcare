package middleware

import (
	"Alice-Seahat-Healthcare/seahat-be/config"
	"Alice-Seahat-Healthcare/seahat-be/constant"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) Cors() gin.HandlerFunc {
	allowOrigins := []string{"http://localhost:3000", "http://127.0.0.1:3000"}

	if config.App.Env == constant.Production {
		allowOrigins = []string{"*"}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "HEAD", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
}
