package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Tortik3000/todo-list/config"
)

func TestNew(t *testing.T) {

	tests := []struct {
		name     string
		envValue string
		setEnv   bool
		wantPort string
	}{
		{
			name:     "success custom port",
			envValue: "9090",
			setEnv:   true,
			wantPort: "9090",
		},
		{
			name:     "success default port when empty",
			envValue: "",
			setEnv:   true,
			wantPort: config.DefaultPort,
		},
		{
			name:     "success default port when unset",
			envValue: "",
			setEnv:   false,
			wantPort: config.DefaultPort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv("PORT", tt.envValue)
			} else {
				os.Unsetenv("PORT")
			}
			
			cfg := config.New()

			assert.Equal(t, tt.wantPort, cfg.Server.Port)
		})
	}
}
