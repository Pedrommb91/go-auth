package services

import (
	"github.com/Pedrommb91/go-auth/config"
	"github.com/Pedrommb91/go-auth/internal/api/models"
	"github.com/Pedrommb91/go-auth/pkg/encrypt"
	"github.com/Pedrommb91/go-auth/pkg/errors"
)

type UserService struct {
	r         models.UserRepositoryInterface
	encrypt   config.Encrypt
	encryptor encrypt.Encryptor
}

type UserServiceInterface interface {
	AddUser(username, email, password string) (int64, error)
}

func NewUserService(r models.UserRepositoryInterface, encrypt config.Encrypt, encryptor encrypt.Encryptor) UserService {
	return UserService{
		r:         r,
		encrypt:   encrypt,
		encryptor: encryptor,
	}
}

func (s UserService) AddUser(username, email, password string) (int64, error) {
	const op errors.Op = "services.AddUser"

	salt := s.encryptor.GenerateSalt(64, true, true)
	passHash, err := s.encryptor.Encrypt(password, salt, s.encrypt.Password)
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
