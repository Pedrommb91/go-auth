package services

import (
	"fmt"
	"testing"

	"github.com/Pedrommb91/go-auth/config"
	"github.com/Pedrommb91/go-auth/internal/api/models"
	"github.com/Pedrommb91/go-auth/mocks"
	"github.com/Pedrommb91/go-auth/pkg/errors"
	"github.com/go-faker/faker/v4"
	uuid "github.com/satori/go.uuid"
)

func TestUserService_AddUser(t *testing.T) {
	dummyID := faker.UUIDHyphenated()
	errors.NewUUID = func() uuid.UUID {
		return uuid.FromStringOrNil(dummyID)
	}

	type addUserMockResponse struct {
		response int64
		err      error
	}
	type encryptMockResponse struct {
		err error
	}
	type fields struct {
		encrypt config.Encrypt
	}
	type args struct {
		username string
		email    string
		password string
		salt     string
	}
	tests := []struct {
		name                string
		addUserMockResponse addUserMockResponse
		encryptMockResponse encryptMockResponse
		fields              fields
		args                args
		want                int64
		expectedErr         error
	}{
		{
			name: "Success",
			addUserMockResponse: addUserMockResponse{
				response: 1,
				err:      nil,
			},
			encryptMockResponse: encryptMockResponse{
				err: nil,
			},
			fields: fields{
				encrypt: config.Encrypt{
					Password: faker.Password(),
				},
			},
			args: args{
				username: faker.Username(),
				email:    faker.Email(),
				password: faker.Password(),
				salt:     faker.Password(),
			},
			want:        1,
			expectedErr: nil,
		},
		{
			name: "Fails to encrypt",
			addUserMockResponse: addUserMockResponse{
				response: 1,
				err:      nil,
			},
			encryptMockResponse: encryptMockResponse{
				err: errors.Build(
					errors.WithError(fmt.Errorf("failed o encrypt password")),
				),
			},
			fields: fields{
				encrypt: config.Encrypt{
					Password: faker.Password(),
				},
			},
			args: args{
				username: faker.Username(),
				email:    faker.Email(),
				password: faker.Password(),
				salt:     faker.Password(),
			},
			want: 0,
			expectedErr: errors.Build(
				errors.WithError(fmt.Errorf("failed o encrypt password")),
			),
		},
		{
			name: "Fails to add user",
			addUserMockResponse: addUserMockResponse{
				response: 0,
				err: errors.Build(
					errors.WithError(fmt.Errorf("failed to add user")),
				),
			},
			encryptMockResponse: encryptMockResponse{
				err: nil,
			},
			fields: fields{
				encrypt: config.Encrypt{
					Password: faker.Password(),
				},
			},
			args: args{
				username: faker.Username(),
				email:    faker.Email(),
				password: faker.Password(),
				salt:     faker.Password(),
			},
			want: 0,
			expectedErr: errors.Build(
				errors.WithError(fmt.Errorf("failed to add user")),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mocks.NewUserRepositoryInterface(t)
			r.On("AddUser", models.Users{
				Username: tt.args.username,
				Email:    tt.args.email,
				Credentials: models.Credentials{
					Salt:     tt.args.salt,
					PassHash: tt.args.password,
				},
			}).Return(tt.addUserMockResponse.response, tt.addUserMockResponse.err).Maybe()

			enc := mocks.NewEncryptor(t)
			enc.On("GenerateSalt", 64, true, true).Return(tt.args.salt).Maybe()
			enc.On("Encrypt", tt.args.password, tt.args.salt, tt.fields.encrypt.Password).Return(tt.args.password, tt.encryptMockResponse.err).Maybe() // no encryption

			s := NewUserService(r, tt.fields.encrypt, enc)
			got, err := s.AddUser(tt.args.username, tt.args.email, tt.args.password)
			if !errors.Equal(errors.GetFirstNestedError(err), tt.expectedErr) {
				t.Errorf("UserService.AddUser() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserService.AddUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
