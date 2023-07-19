package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/Pedrommb91/go-auth/internal/api/models"
	"github.com/Pedrommb91/go-auth/pkg/errors"
	"github.com/go-faker/faker/v4"
	"github.com/rs/zerolog"
)

func (s *DatabaseTestSuite) TestCreateWithRelation() {
	now := time.Unix(faker.UnixTime(), 0).UTC()
	type args struct {
		user models.Users
	}
	tests := []struct {
		name    string
		agrs    args
		wantId  int64
		wantErr error
	}{
		{
			name: "Insert John with success",
			agrs: args{
				user: models.Users{
					Username: "John",
					Email:    faker.Email(),
					Credentials: models.Credentials{
						Salt:     faker.Password(),
						PassHash: faker.Password(),
					},
					CreatedAt: now,
				},
			},
			wantId:  1,
			wantErr: nil,
		},
		{
			name: "Insert John again will throw an error as username must be unique",
			agrs: args{
				user: models.Users{
					Username: "John",
					Email:    faker.Email(),
					Credentials: models.Credentials{
						Salt:     faker.Password(),
						PassHash: faker.Password(),
					},
					CreatedAt: now,
				},
			},
			wantId: 0,
			wantErr: errors.Build(
				errors.WithOp("database.createWithParentRelations"),
				errors.WithError(fmt.Errorf("pq: duplicate key value violates unique constraint \"users_username_key\"")),
				errors.WithMessage("Constrain violation: failed to insert entry"),
				errors.KindBadRequest(),
				errors.WithSeverity(zerolog.WarnLevel),
			),
		},
		{
			name: "Insert without credentials",
			agrs: args{
				user: models.Users{
					Username:  faker.Username(),
					Email:     faker.Email(),
					CreatedAt: now,
				},
			},
			wantId: 0,
			wantErr: errors.Build(
				errors.WithOp("database.createWithParentRelations"),
				errors.WithError(fmt.Errorf("pq: null value in column \"salt\" of relation \"credentials\" violates not-null constraint")),
				errors.WithMessage("Constrain violation: failed to insert entry"),
				errors.KindBadRequest(),
				errors.WithSeverity(zerolog.WarnLevel),
			),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := With[models.Users](s.db).Insert(tt.agrs.user)
			if !errors.Equal(errors.GetFirstNestedError(err), tt.wantErr) {
				s.T().Errorf("Insert() = %v, want %v", err, tt.wantErr)
			}
			if got != tt.wantId {
				s.T().Errorf("Insert() = %v, want %v", got, tt.wantId)
			}
		})
	}
}

func (s *DatabaseTestSuite) TestCreateWithoutRelation() {
	type args struct {
		credentials models.Credentials
	}
	tests := []struct {
		name    string
		agrs    args
		wantId  int64
		wantErr error
	}{
		{
			name: "Success",
			agrs: args{
				credentials: models.Credentials{
					Salt:     faker.Password(),
					PassHash: faker.Password(),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := With[models.Credentials](s.db).Insert(tt.agrs.credentials)
			if !errors.Equal(errors.GetFirstNestedError(err), tt.wantErr) {
				s.T().Errorf("Insert() = %v, want %v", err, tt.wantErr)
			}
			if got == 0 {
				s.T().Errorf("Insert() = %v, want != 0", got)
			}
		})
	}
}
