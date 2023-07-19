package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Pedrommb91/go-auth/pkg/errors"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type DatabaseTestSuite struct {
	container *ContainerDBConfigs
	suite.Suite
	db *sql.DB
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (s *DatabaseTestSuite) SetupSuite() {
	dummyID := "e157f89f-abd0-4b1a-bc58-de8bd8fd04cd"
	errors.NewUUID = func() uuid.UUID {
		return uuid.FromStringOrNil(dummyID)
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	s.container = NewPostgresTestContainer(ctx, "RepositoryPostgres")

	var err error
	s.db = NewPostgresOrDie(s.container.Config)

	err = s.container.RunMigrations()
	if err != nil {
		s.Require().NoError(err)
	}
}

func (s *DatabaseTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	s.Require().NoError(s.container.Container.Terminate(ctx))
}

func TestNewPostgresOrDie(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Conection or die",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewPostgresTestContainer(context.Background(), "ConnectionDb")

			db := NewPostgresOrDie(container.Config)
			if db == nil {
				t.Errorf("Failed to connect to database")
			}
		})
	}
}
