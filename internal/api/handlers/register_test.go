package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Pedrommb91/go-auth/config"
	"github.com/Pedrommb91/go-auth/internal/api/middlewares"
	"github.com/Pedrommb91/go-auth/internal/api/openapi"
	"github.com/Pedrommb91/go-auth/mocks"
	"github.com/Pedrommb91/go-auth/pkg/errors"
	"github.com/Pedrommb91/go-auth/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_client_RegisterUserHandler(t *testing.T) {
	dummyID := faker.UUIDHyphenated()
	errors.NewUUID = func() uuid.UUID {
		return uuid.FromStringOrNil(dummyID)
	}

	now := time.Unix(faker.UnixTime(), 0).UTC()

	path := "/api/v1/register"

	type addUserMockResponse struct {
		response int64
		err      error
	}
	type fields struct {
		cfg *config.Config
		log logger.Interface
	}
	type args struct {
		requestBody *openapi.RegisterUserRequestBody
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		addUserMockResponse   addUserMockResponse
		expectedErrorResponse *openapi.Error
		expectedCode          int
	}{
		{
			name: "Success",
			fields: fields{
				cfg: &config.Config{},
				log: logger.New("info"),
			},
			args: args{
				requestBody: &openapi.RegisterUserRequestBody{
					Email:    faker.Email(),
					Password: "#sdjU1kaL!",
					Username: faker.Username(),
				},
			},
			addUserMockResponse: addUserMockResponse{
				response: 1,
				err:      nil,
			},
			expectedErrorResponse: nil,
			expectedCode:          http.StatusCreated,
		},
		{
			name: "Invalid user request body",
			fields: fields{
				cfg: &config.Config{},
				log: logger.New("info"),
			},
			args: args{
				requestBody: &openapi.RegisterUserRequestBody{
					Email:    faker.Email(),
					Password: "#sdjU1kaL!",
					Username: "a",
				},
			},
			addUserMockResponse: addUserMockResponse{
				response: 1,
				err:      nil,
			},
			expectedErrorResponse: &openapi.Error{
				Error:     "Bad Request",
				Id:        dummyID,
				Message:   "The size of the username must be more than 3 and less than 64",
				Path:      path,
				Status:    http.StatusBadRequest,
				Timestamp: now,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Empty request body",
			fields: fields{
				cfg: &config.Config{},
				log: logger.New("info"),
			},
			args: args{
				requestBody: nil,
			},
			addUserMockResponse: addUserMockResponse{
				response: 1,
				err:      nil,
			},
			expectedErrorResponse: &openapi.Error{
				Error:     "Bad Request",
				Id:        dummyID,
				Message:   "Invalid user",
				Path:      path,
				Status:    http.StatusBadRequest,
				Timestamp: now,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Error adding user to db",
			fields: fields{
				cfg: &config.Config{},
				log: logger.New("info"),
			},
			args: args{
				requestBody: &openapi.RegisterUserRequestBody{
					Email:    faker.Email(),
					Password: "#sdjU1kaL!",
					Username: faker.Username(),
				},
			},
			addUserMockResponse: addUserMockResponse{
				response: 0,
				err: errors.Build(
					errors.WithError(fmt.Errorf("failed to add user")),
					errors.WithMessage("Failed to add user"),
				),
			},
			expectedErrorResponse: &openapi.Error{
				Error:     "Unexpected Error",
				Id:        dummyID,
				Message:   "Failed to add user",
				Path:      path,
				Status:    http.StatusInternalServerError,
				Timestamp: now,
			},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()

			userServiceMock := mocks.NewUserServiceInterface(t)
			if tt.args.requestBody != nil {
				userServiceMock.On(
					"AddUser",
					tt.args.requestBody.Username,
					tt.args.requestBody.Email,
					tt.args.requestBody.Password).
					Return(tt.addUserMockResponse.response, tt.addUserMockResponse.err).Maybe()
			}

			services := &Services{
				User: userServiceMock,
			}

			clockMock := mocks.NewClock(t)
			clockMock.On("Now").Return(now).Maybe()

			r.Use(middlewares.ErrorHandler(clockMock, tt.fields.log))

			g := NewClient(tt.fields.cfg, tt.fields.log, services)
			r.POST(path, func(c *gin.Context) {
				g.RegisterUserHandler(c)
			})

			w := httptest.NewRecorder()
			var req *http.Request
			if tt.args.requestBody != nil {
				data, err := json.Marshal(tt.args.requestBody)
				if err != nil {
					t.Errorf("Failed to marshal request body")
				}
				req, _ = http.NewRequest(http.MethodPost, path, bytes.NewReader(data))
			} else {
				req, _ = http.NewRequest(http.MethodPost, path, nil)
			}
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode != http.StatusCreated {
				var got *openapi.Error
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Errorf("Failed to unmarshal body: %s", err)
				}
				assert.Equal(t, tt.expectedErrorResponse, got)
			}
		})
	}
}
