package repositories

import (
	"database/sql"

	"github.com/Pedrommb91/go-auth/internal/api/models"
	"github.com/Pedrommb91/go-auth/pkg/database"
	"github.com/Pedrommb91/go-auth/pkg/errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) AddUser(user models.Users) (int64, error) {
	const op errors.Op = "repositories.AddUser"

	id, err := database.With[models.Users](r.db).Create(user)
	if err != nil {
		return 0, errors.Build(
			errors.WithOp(op),
			errors.WithError(err),
			errors.WithMessage("Failed to register user"),
		)
	}

	return id, nil
}
