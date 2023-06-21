package services

import (
	"github.com/Pedrommb91/go-auth/config"
	"github.com/Pedrommb91/go-auth/internal/api/models"
	"github.com/Pedrommb91/go-auth/pkg/encrypt"
	"github.com/Pedrommb91/go-auth/pkg/errors"
)

type UserService struct {
	r       models.UserRepositoryInterface
	encrypt config.Encrypt
}

type UserServiceInterface interface {
	AddUser(username, email, password string) (int64, error)
}

func NewUserService(r models.UserRepositoryInterface, encrypt config.Encrypt) UserService {
	return UserService{
		r:       r,
		encrypt: encrypt,
	}
}

func (s UserService) AddUser(username, email, password string) (int64, error) {
	const op errors.Op = "services.AddUser"

	salt := encrypt.GeneratePassword(64, true, true)
	e := encrypt.PasswordEncryptor{}
	passHash, err := e.Encrypt(password, salt, s.encrypt.Password)
	if err != nil {
		return 0, errors.Build(
			errors.WithOp(op),
			errors.WithError(err),
			errors.WithMessage("Failed to enryp password"),
		)
	}

	id, err := s.r.AddUser(models.Users{
		Username: username,
		Email:    email,
		Credentials: models.Credentials{
			Salt:     salt,
			PassHash: passHash,
		},
	})
	if err != nil {
		return 0, errors.Build(
			errors.WithOp(op),
			errors.WithError(err),
			errors.WithMessage("Failed to register user"),
		)
	}

	return id, nil
}
