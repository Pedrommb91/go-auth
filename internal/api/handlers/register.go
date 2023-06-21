package handlers

import (
	"net/http"

	"github.com/Pedrommb91/go-auth/internal/api/openapi"
	"github.com/Pedrommb91/go-auth/pkg/errors"
	"github.com/gin-gonic/gin"
)

// RegisterUserHandler implements openapi.ServerInterface.
func (cli *client) RegisterUserHandler(c *gin.Context) {
	const op errors.Op = "handlers.RegisterUserHandler"

	var user *openapi.RegisterUserRequestBody
	if err := c.ShouldBind(&user); err != nil {
		c.Error(errors.Build(
			errors.WithOp(op),
			errors.WithError(err),
			errors.WithMessage("Invalid user"),
			errors.KindBadRequest(),
		))
		return
	}
	id, err := cli.services.User.AddUser(user.Username, user.Email, user.Password)
	if err != nil {
		c.Error(errors.Build(
			errors.WithOp(op),
			errors.WithError(err),
			errors.WithMessage("Failed to register user"),
		))
		return
	}

	c.JSON(http.StatusCreated, &openapi.CreateUserResponse{
		Id: id,
	})
}
