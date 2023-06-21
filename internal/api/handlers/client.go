package handlers

import (
	"github.com/Pedrommb91/go-auth/config"
	"github.com/Pedrommb91/go-auth/internal/api/openapi"
	"github.com/Pedrommb91/go-auth/internal/api/services"
	"github.com/Pedrommb91/go-auth/pkg/logger"
)

type client struct {
	cfg      *config.Config
	log      logger.Interface
	services *Services
}

type Services struct {
	User services.UserServiceInterface
}

func NewClient(cfg *config.Config, l logger.Interface, services *Services) openapi.ServerInterface {
	return &client{
		cfg:      cfg,
		log:      l,
		services: services,
	}
}
