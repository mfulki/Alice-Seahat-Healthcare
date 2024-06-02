package utils

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/config"
	"Alice-Seahat-Healthcare/seahat-be/constant"

	"github.com/gin-gonic/gin"
)

func SetCookie(ctx *gin.Context, name string, value string, duration time.Duration) {
	secure := false
	if config.App.Env == "production" {
		secure = true
	}

	ctx.SetCookie(name, value, int(duration.Seconds()), "/", DomainURL(config.FE.URL), secure, false)
}

func SendLoginCookie(ctx *gin.Context, token string) {
	SetCookie(ctx, constant.TokenAccess, token, constant.JwtDefaultDuration)
}

func SendLogoutCookie(ctx *gin.Context) {
	SetCookie(ctx, constant.TokenAccess, "", 0)
}
