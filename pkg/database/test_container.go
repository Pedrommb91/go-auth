package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/Pedrommb91/go-auth/config"
	"github.com/docker/go-connections/nat"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ContainerDBConfigs struct {
	Config    config.Database
	Container testcontainers.Container
	DSN       string
}

func NewPostgresTestContainer(ctx context.Context, name string) *ContainerDBConfigs {
	cfg := config.Database{
		Host:     "localhost",
		User:     "user",
		Password: "password",
		DBName:   "db",
		Port:     "5432/tcp",
		SslMode:  "disable",
		Schema:   "public",
	}

	req := testcontainers.ContainerRequest{
		Name: name,
		Env: map[string]string{
			"POSTGRES_USER":     cfg.User,
			"POSTGRES_PASSWORD": cfg.Password,
			"POSTGRES_DB":       cfg.DBName,
		},
		ExposedPorts: []string{cfg.Port},
		Image:        "postgres:14.3",
		WaitingFor: wait.ForAll(
			wait.ForExec([]string{"pg_isready"}).
				WithPollInterval(2 * time.Second).
				WithExitCodeMatcher(func(exitCode int) bool {
					return exitCode == 0
				}),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		panic(err)
	}

	cfg.Host, err = container.Host(ctx)
	if err != nil {
		panic(err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(cfg.Port))
	if err != nil {
		panic(err)
	}

	cfg.Port = mappedPort.Port()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SslMode)

	return &ContainerDBConfigs{
		Config:    cfg,
		Container: container,
		DSN:       dsn,
	}
}

func (c ContainerDBConfigs) RunMigrations() error {
	var sqlMigrations *sql.DB
	sqlMigrations, err := sql.Open("postgres", c.DSN)
	if err != nil {
		return err
	}
	// nolint: dogsled
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../migrations")

	files := os.DirFS(dir)
	goose.SetBaseFS(files)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(sqlMigrations, "."); err != nil {
		panic(err)
	}
	return nil
}
