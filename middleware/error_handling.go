package middleware

import (
	"errors"
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	cv "Alice-Seahat-Healthcare/seahat-be/libs/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func (m *Middleware) ErrorHandling(log *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		firstError := ctx.Errors[0]

		if appError := new(apperror.AppError); errors.As(firstError, &appError) {
			ctx.AbortWithStatusJSON(appError.StatusCode, response.Body{
				Message: appError.Error(),
			})

			return
		}

		if ve := make(validator.ValidationErrors, 0); errors.As(firstError, &ve) {
			mapErrors := make(map[string]string)
			for _, fe := range ve {
				mapErrors[fe.Field()] = cv.GetValidationErrorMsg(fe)
			}

			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Body{
				Message: apperror.ErrInvalidRequest.Error(),
				Errors:  mapErrors,
			})

			return
		}

		log.Error(firstError)

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Body{
			Message: apperror.ErrInternalServer.Error(),
		})
	}
}
