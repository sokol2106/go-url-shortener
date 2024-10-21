package test

import (
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	testData := struct {
		ServerAddress   string
		BaseURL         string
		FileStoragePath string
		DatabaseDsn     string
		EnableHTTPS     string
	}{
		ServerAddress:   "127.0.0.1:8080",
		BaseURL:         "http://localhost:9090",
		FileStoragePath: "./test/",
		DatabaseDsn:     "connect",
		EnableHTTPS:     "true",
	}

	t.Run("testConfig", func(t *testing.T) {
		cnf := config.NewConfigURL(
			testData.ServerAddress,
			"http://localhost:9090",
			testData.FileStoragePath,
			testData.DatabaseDsn,
			testData.EnableHTTPS,
		)

		assert.Equal(t, testData.ServerAddress, cnf.ServerAddress())
		assert.Equal(t, testData.BaseURL, cnf.BaseURL())
		assert.Equal(t, testData.FileStoragePath, cnf.FileStoragePath())
		assert.Equal(t, testData.DatabaseDsn, cnf.DatabaseDsn())
		assert.Equal(t, true, cnf.EnableHTTPS())

		cnf.SetEnableHTTPS("")

		assert.Equal(t, false, cnf.EnableHTTPS())

	})
}
