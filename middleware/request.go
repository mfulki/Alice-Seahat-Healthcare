package middleware

import (
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) RequestID(c *gin.Context) {
	randomStr, err := utils.RandomString(constant.LengthOfRequestID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Set(constant.RequestID, randomStr)
	c.Next()
}
