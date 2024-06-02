package middleware

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/constant"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (m *Middleware) Logger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		statusCode := c.Writer.Status()
		requestID, exist := c.Get(constant.RequestID)
		if !exist {
			requestID = ""
		}

		entry := log.WithFields(logrus.Fields{
			"path":        path,
			"method":      c.Request.Method,
			"latency":     time.Since(start),
			"request_id":  requestID,
			"status_code": statusCode,
		})

		if statusCode >= 500 && statusCode <= 599 {
			entry.WithField("error", c.Errors[0]).Error("encountered error")
			return
		}

		entry.Info("request processed")
	}
}
