package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("GIN_MODE", "test")
	os.Setenv("MONGODB_URI", "mongodb://test:27017")
	defer os.Unsetenv("SERVER_PORT")
	defer os.Unsetenv("GIN_MODE")
	defer os.Unsetenv("MONGODB_URI")

	cfg, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "9090", cfg.Server.Port)
	assert.Equal(t, "test", cfg.Server.Mode)
	assert.Equal(t, "mongodb://test:27017", cfg.Database.URI)
}
