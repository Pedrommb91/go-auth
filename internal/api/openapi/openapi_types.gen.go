// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package openapi

import (
	"time"
)

// CreateUserResponse defines model for CreateUserResponse.
type CreateUserResponse struct {
	Id int64 `json:"id"`
}

// Error defines model for Error.
type Error struct {
	Error     string    `json:"error"`
	Id        string    `json:"id"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
	Status    int32     `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// RegisterUserRequestBody defines model for RegisterUserRequestBody.
type RegisterUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// RegisterUserHandlerJSONRequestBody defines body for RegisterUserHandler for application/json ContentType.
type RegisterUserHandlerJSONRequestBody = RegisterUserRequestBody
