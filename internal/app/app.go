package app

import (
	"database/sql"

	"github.com/Pedrommb91/go-auth/config"
	"github.com/Pedrommb91/go-auth/internal/api"
	"github.com/Pedrommb91/go-auth/internal/api/handlers"
	"github.com/Pedrommb91/go-auth/internal/api/repositories"
	"github.com/Pedrommb91/go-auth/internal/api/services"
	"github.com/Pedrommb91/go-auth/pkg/database"
	"github.com/Pedrommb91/go-auth/pkg/logger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	db := database.NewPostgresOrDie(cfg.Database)
	services := createServices(db, cfg)

	server := api.NewServer(cfg, l)
	server.ServerConfigure()
	server.SetRoutes(services)
	server.Run()
}

func createServices(db *sql.DB, cfg *config.Config) *handlers.Services {
	ur := repositories.NewUserRepository(db)
	return &handlers.Services{
		User: services.NewUserService(ur, cfg.Encrypt),
	}
}
