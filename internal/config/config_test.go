package config_test

import (
	"testing"

	"github.com/Financial-Partner/server/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDataDir = "testdata"

func TestLoadConfig(t *testing.T) {
	t.Run("Valid YAML config", func(t *testing.T) {
		cfg, err := config.LoadConfig(testDataDir + "/valid")
		require.NoError(t, err)

		assert.Equal(t, "mongodb://user:pass@localhost:27017", cfg.MongoDB.URI)
		assert.Equal(t, "testdb", cfg.MongoDB.Database)
		assert.Equal(t, "localhost:6379", cfg.Redis.Host)
		assert.Equal(t, "redispass", cfg.Redis.Password)
		assert.Equal(t, 1, cfg.Redis.DB)
		assert.Equal(t, "test-project", cfg.Firebase.ProjectID)
		assert.Equal(t, "creds.json", cfg.Firebase.CredentialFile)
	})

	t.Run("Invalid YAML format", func(t *testing.T) {
		_, err := config.LoadConfig(testDataDir + "/invalid")
		assert.Error(t, err)
	})

	t.Run("Only environment variables", func(t *testing.T) {
		t.Setenv("PARTNER_SERVER_MONGODB_URI", "mongodb://env:pass@localhost:27017")
		t.Setenv("PARTNER_SERVER_MONGODB_DATABASE", "envdb")
		t.Setenv("PARTNER_SERVER_REDIS_HOST", "env-redis:6379")
		t.Setenv("PARTNER_SERVER_REDIS_PASSWORD", "env-redispass")
		t.Setenv("PARTNER_SERVER_REDIS_DB", "2")
		t.Setenv("PARTNER_SERVER_FIREBASE_PROJECT_ID", "env-project")
		t.Setenv("PARTNER_SERVER_FIREBASE_CREDENTIAL_FILE", "env-creds.json")

		cfg, err := config.LoadConfig(testDataDir + "/not_found")
		require.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("Environment variables override YAML", func(t *testing.T) {
		t.Setenv("PARTNER_SERVER_MONGODB_URI", "mongodb://override:pass@localhost:27017")
		t.Setenv("PARTNER_SERVER_REDIS_PASSWORD", "override-redispass")
		t.Setenv("PARTNER_SERVER_FIREBASE_PROJECT_ID", "override-project")

		cfg, err := config.LoadConfig(testDataDir + "/valid")
		require.NoError(t, err)

		assert.Equal(t, "mongodb://override:pass@localhost:27017", cfg.MongoDB.URI)
		assert.Equal(t, "testdb", cfg.MongoDB.Database)
		assert.Equal(t, "localhost:6379", cfg.Redis.Host)
		assert.Equal(t, "override-redispass", cfg.Redis.Password)
		assert.Equal(t, 1, cfg.Redis.DB)
		assert.Equal(t, "override-project", cfg.Firebase.ProjectID)
		assert.Equal(t, "creds.json", cfg.Firebase.CredentialFile)
	})

	t.Run("Missing fields in YAML", func(t *testing.T) {
		cfg, err := config.LoadConfig(testDataDir + "/incomplete")
		require.NoError(t, err)

		assert.Equal(t, "mongodb://incomplete:pass@localhost:27017", cfg.MongoDB.URI)
		assert.Equal(t, "", cfg.MongoDB.Database)
		assert.Equal(t, "incomplete-redis:6379", cfg.Redis.Host)
		assert.Equal(t, "", cfg.Redis.Password)
		assert.Equal(t, 0, cfg.Redis.DB)
		assert.Equal(t, "", cfg.Firebase.ProjectID)
		assert.Equal(t, "", cfg.Firebase.CredentialFile)
	})
}
