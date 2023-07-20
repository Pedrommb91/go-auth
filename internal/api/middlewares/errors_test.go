package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Pedrommb91/go-auth/internal/api/openapi"
	"github.com/Pedrommb91/go-auth/mocks"
	"github.com/Pedrommb91/go-auth/pkg/errors"
	"github.com/Pedrommb91/go-auth/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	path := "/"

	dummyID := faker.UUIDHyphenated()
	errors.NewUUID = func() uuid.UUID {
		return uuid.FromStringOrNil(dummyID)
	}

	now := time.Unix(faker.UnixTime(), 0).UTC()

	type args struct {
		err error
	}
	tests := []struct {
		name                  string
		args                  args
		noErrors              bool
		expectedErrorResponse *openapi.Error
		expectedCode          int
	}{
		{
			name: "Error as response",
			args: args{
				err: errors.Build(
					errors.WithError(fmt.Errorf("dummy error")),
					errors.WithMessage("Dummy error"),
				),
			},
			noErrors: false,
			expectedErrorResponse: &openapi.Error{
				Error:     "Unexpected Error",
				Id:        dummyID,
				Message:   "Dummy error",
				Path:      path,
				Status:    http.StatusInternalServerError,
				Timestamp: now,
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "No errors",
			args: args{
				err: nil,
			},
			noErrors:              true,
			expectedErrorResponse: nil,
			expectedCode:          http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clockMock := mocks.NewClock(t)
			clockMock.On("Now").Return(now).Maybe()

			r := gin.Default()
			r.Use(ErrorHandler(clockMock, logger.New("info")))
			r.GET(path, func(c *gin.Context) {
				if !tt.noErrors {
					c.Error(fmt.Errorf("ignore this type of errors"))
					c.Error(tt.args.err)
				}
			})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, path, nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode != http.StatusOK {
				var got *openapi.Error
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Errorf("Failed to unmarshal body: %s", err)
				}
				assert.Equal(t, tt.expectedErrorResponse, got)
			}

		})
	}
}
