package handlers

import (
	"github.com/Pedrommb91/go-auth/config"
	"github.com/Pedrommb91/go-auth/internal/api/openapi"
	"github.com/Pedrommb91/go-auth/internal/api/repositories"
	"github.com/Pedrommb91/go-auth/pkg/logger"
)

type client struct {
	cfg *config.Config
	log logger.Interface
	db  *repositories.PostgresDB
}

func NewClient(cfg *config.Config, l logger.Interface, db *repositories.PostgresDB) openapi.ServerInterface {
	return &client{
		cfg: cfg,
		log: l,
		db:  db,
	}
}
