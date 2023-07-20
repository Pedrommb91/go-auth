package models

import (
	"fmt"
	"net/mail"
	"unicode"

	"github.com/Pedrommb91/go-auth/internal/api/openapi"
	"github.com/Pedrommb91/go-auth/pkg/errors"
	"github.com/rs/zerolog"
)

type RegisterUserRequestBody openapi.RegisterUserRequestBody

func (b RegisterUserRequestBody) Validate() error {
	const op errors.Op = "models.Validate"
	if len(b.Username) < 3 || len(b.Username) > 63 {
		return errors.Build(
			errors.WithOp(op),
			errors.WithError(fmt.Errorf("the length of the username should be more than 3 and less than 64")),
			errors.WithMessage("The size of the username must be more than 3 and less than 64"),
			errors.KindBadRequest(),
			errors.WithSeverity(zerolog.WarnLevel),
		)
	}

	if err := b.verifyPassword(); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(b.Email); err != nil {
		return errors.Build(
			errors.WithOp(op),
			errors.WithError(err),
			errors.WithMessage("Invalid email"),
			errors.KindBadRequest(),
			errors.WithSeverity(zerolog.WarnLevel),
		)
	}

	return nil
}

func (b RegisterUserRequestBody) verifyPassword() error {
	const op errors.Op = "models.verifyPassword"
	if len(b.Password) < 8 || len(b.Password) > 16 {
		return errors.Build(
			errors.WithOp(op),
			errors.WithError(fmt.Errorf("password must have a minimum of eight characters and a maximum of 16")),
			errors.WithMessage("Password must have a minimum of eight characters and a maximum of 16"),
			errors.KindBadRequest(),
			errors.WithSeverity(zerolog.WarnLevel),
		)
	}

	var number, upper, lower, special bool
	for _, c := range b.Password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsSpace(c):
			return errors.Build(
				errors.WithOp(op),
				errors.WithError(fmt.Errorf("spaces in password is not allowed")),
				errors.WithMessage("Spaces in password is not allowed"),
				errors.KindBadRequest(),
				errors.WithSeverity(zerolog.WarnLevel),
			)
		}
	}

	if !number || !upper || !lower || !special {
		return errors.Build(
			errors.WithOp(op),
			errors.WithError(fmt.Errorf("password must have one upper case, one lower case, one number and a special char")),
			errors.WithMessage("Password must have one upper case, one lower case, one number and a special char"),
			errors.KindBadRequest(),
			errors.WithSeverity(zerolog.WarnLevel),
		)
	}

	return nil
}
