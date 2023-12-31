package router

import (
	"testing"

	"github.com/Pedrommb91/go-auth/config"
	"github.com/Pedrommb91/go-auth/pkg/logger"
	"github.com/gin-gonic/gin"
)

func TestNewRouter(t *testing.T) {
	type args struct {
		engine *gin.Engine
		l      logger.Interface
		cfg    *config.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "New router success",
			args: args{
				engine: &gin.Engine{RouterGroup: gin.New().RouterGroup},
				l:      logger.New("info"),
				cfg:    &config.Config{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewRouter(tt.args.engine, tt.args.l, tt.args.cfg, nil)
			if len(tt.args.engine.Handlers) == 0 {
				t.Errorf("Failed to register handlers")
			}
		})
	}
}
