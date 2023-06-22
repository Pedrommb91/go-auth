package middlewares

import (
	"github.com/Pedrommb91/go-auth/internal/api/openapi"
	"github.com/Pedrommb91/go-auth/pkg/clock"
	"github.com/Pedrommb91/go-auth/pkg/errors"
	"github.com/Pedrommb91/go-auth/pkg/logger"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(clock clock.Clock, l logger.Interface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		for _, v := range ctx.Errors {
			err, ok := v.Err.(*errors.Error)
			if !ok {
				l.Error("Unexpected error: %s", v.Err)
				continue
			}
			er := errors.GetFirstNestedError(err)
			err, _ = er.(*errors.Error)

			ctx.JSON(err.Kind.Int(), &openapi.Error{
				Error:     err.Kind.String(),
				Id:        err.ID,
				Message:   err.Message,
				Path:      ctx.FullPath(),
				Status:    int32(err.Kind.Int()),
				Timestamp: clock.Now(),
			})
			return
		}
	}
}
