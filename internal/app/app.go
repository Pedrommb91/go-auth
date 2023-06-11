package app

import (
	"github.com/Pedrommb91/go-auth/config"
	"github.com/Pedrommb91/go-auth/internal/api"
	"github.com/Pedrommb91/go-auth/internal/api/repositories"
	"github.com/Pedrommb91/go-auth/pkg/database"
	"github.com/Pedrommb91/go-auth/pkg/logger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	db := &repositories.PostgresDB{DB: database.NewPostgresOrDie(cfg.Database)}

	server := api.NewServer(cfg, l, db)
	server.ServerConfigure()
	server.SetRoutes()
	server.Run()
}
