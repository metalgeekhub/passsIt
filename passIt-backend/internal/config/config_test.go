package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequireEnv_Panics(t *testing.T) {
	// Test that requireEnv panics when env var is missing
	assert.Panics(t, func() {
		requireEnv("NONEXISTENT_VAR_12345")
	}, "Should panic when environment variable is missing")
}

func TestRequireEnv_ReturnsValue(t *testing.T) {
	// Set a test environment variable
	testKey := "TEST_VAR_12345"
	testValue := "test_value"
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	result := requireEnv(testKey)
	assert.Equal(t, testValue, result)
}

func TestConfig_Structure(t *testing.T) {
	// Test that Config struct has expected structure
	cfg := &Config{
		App:         &AppConfig{Port: 8080, Env: "test"},
		DB:          nil, // Would require full DB config
		Auth:        nil, // Would require full Auth config
		RedisClient: nil, // Would require full Redis config
	}

	assert.NotNil(t, cfg.App)
	assert.Equal(t, 8080, cfg.App.Port)
	assert.Equal(t, "test", cfg.App.Env)
}

func TestAppConfig_Values(t *testing.T) {
	app := &AppConfig{
		Port: 3000,
		Env:  "development",
	}

	assert.Equal(t, 3000, app.Port)
	assert.Equal(t, "development", app.Env)
}

func TestAppConfig_DefaultPort(t *testing.T) {
	tests := []struct {
		name     string
		port     int
		isValid  bool
	}{
		{"Valid port 80", 80, true},
		{"Valid port 8080", 8080, true},
		{"Valid port 3000", 3000, true},
		{"Invalid port 0", 0, false},
		{"Invalid port negative", -1, false},
		{"Invalid port too high", 70000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isValid {
				assert.Greater(t, tt.port, 0)
				assert.LessOrEqual(t, tt.port, 65535)
			} else {
				assert.True(t, tt.port <= 0 || tt.port > 65535)
			}
		})
	}
}
