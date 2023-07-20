package models

import (
	"fmt"
	"testing"

	"github.com/Pedrommb91/go-auth/pkg/errors"
	"github.com/go-faker/faker/v4"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
)

func TestRegisterUserRequestBody_Validate(t *testing.T) {
	const op errors.Op = "models.Validate"
	dummyID := faker.UUIDHyphenated()
	errors.NewUUID = func() uuid.UUID {
		return uuid.FromStringOrNil(dummyID)
	}

	tests := []struct {
		name        string
		b           RegisterUserRequestBody
		expectedErr error
	}{
		{
			name: "Success",
			b: RegisterUserRequestBody{
				Email:    faker.Email(),
				Password: "#sdjU1kaL!",
				Username: faker.Username(),
			},
			expectedErr: nil,
		},
		{
			name: "Short length of username",
			b: RegisterUserRequestBody{
				Email:    faker.Email(),
				Password: "#sdjU1kaL!",
				Username: "aa",
			},
			expectedErr: errors.Build(
				errors.WithOp(op),
				errors.WithError(fmt.Errorf("the length of the username should be more than 3 and less than 64")),
				errors.WithMessage("The size of the username must be more than 3 and less than 64"),
				errors.KindBadRequest(),
				errors.WithSeverity(zerolog.WarnLevel),
			),
		},
		{
			name: "Long length of username",
			b: RegisterUserRequestBody{
				Email:    faker.Email(),
				Password: "#sdjU1kaL!",
				Username: faker.Username() + faker.UUIDHyphenated() + faker.UUIDHyphenated(),
			},
			expectedErr: errors.Build(
				errors.WithOp(op),
				errors.WithError(fmt.Errorf("the length of the username should be more than 3 and less than 64")),
				errors.WithMessage("The size of the username must be more than 3 and less than 64"),
				errors.KindBadRequest(),
				errors.WithSeverity(zerolog.WarnLevel),
			),
		},
		{
			name: "Password without special chars",
			b: RegisterUserRequestBody{
				Email:    faker.Email(),
				Password: "sdjU1kaL",
				Username: faker.Username(),
			},
			expectedErr: errors.Build(
				errors.WithOp("models.verifyPassword"),
				errors.WithError(fmt.Errorf("password must have one upper case, one lower case, one number and a special char")),
				errors.WithMessage("Password must have one upper case, one lower case, one number and a special char"),
				errors.KindBadRequest(),
				errors.WithSeverity(zerolog.WarnLevel),
			),
		},
		{
			name: "Password with short length",
			b: RegisterUserRequestBody{
				Email:    faker.Email(),
				Password: "U!1kaL",
				Username: faker.Username(),
			},
			expectedErr: errors.Build(
				errors.WithOp("models.verifyPassword"),
				errors.WithError(fmt.Errorf("password must have a minimum of eight characters and a maximum of 16")),
				errors.WithMessage("Password must have a minimum of eight characters and a maximum of 16"),
				errors.KindBadRequest(),
				errors.WithSeverity(zerolog.WarnLevel),
			),
		},
		{
			name: "Password with long length",
			b: RegisterUserRequestBody{
				Email:    faker.Email(),
				Password: faker.UUIDHyphenated() + "U!1kaL",
				Username: faker.Username(),
			},
			expectedErr: errors.Build(
				errors.WithOp("models.verifyPassword"),
				errors.WithError(fmt.Errorf("password must have a minimum of eight characters and a maximum of 16")),
				errors.WithMessage("Password must have a minimum of eight characters and a maximum of 16"),
				errors.KindBadRequest(),
				errors.WithSeverity(zerolog.WarnLevel),
			),
		},
		{
			name: "Password with spaces",
			b: RegisterUserRequestBody{
				Email:    faker.Email(),
				Password: "#sdjU 1kaL!",
				Username: faker.Username(),
			},
			expectedErr: errors.Build(
				errors.WithOp("models.verifyPassword"),
				errors.WithError(fmt.Errorf("spaces in password is not allowed")),
				errors.WithMessage("Spaces in password is not allowed"),
				errors.KindBadRequest(),
				errors.WithSeverity(zerolog.WarnLevel),
			),
		},
		{
			name: "Invalid email",
			b: RegisterUserRequestBody{
				Email:    faker.Username(),
				Password: "#sdjU1kaL!",
				Username: faker.Username(),
			},
			expectedErr: errors.Build(
				errors.WithOp(op),
				errors.WithError(fmt.Errorf("mail: missing '@' or angle-addr")),
				errors.WithMessage("Invalid email"),
				errors.KindBadRequest(),
				errors.WithSeverity(zerolog.WarnLevel),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.b.Validate()
			if !errors.Equal(errors.GetFirstNestedError(err), tt.expectedErr) {
				t.Errorf("RegisterUserRequestBody.Validate() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}
