package test

import (
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	testData := struct {
		ServerAddress   string
		BaseUrl         string
		FileStoragePath string
		DatabaseDsn     string
		EnableHttps     string
	}{
		ServerAddress:   "127.0.0.1:8080",
		BaseUrl:         "http://127.0.0.1",
		FileStoragePath: "./test/",
		DatabaseDsn:     "connect",
		EnableHttps:     "true",
	}

	t.Run("testConfig", func(t *testing.T) {
		cnf := config.NewConfigURL(
			testData.ServerAddress,
			testData.BaseUrl,
			testData.FileStoragePath,
			testData.DatabaseDsn,
			testData.EnableHttps,
		)

		assert.Equal(t, testData.ServerAddress, cnf.ServerAddress())
		assert.Equal(t, testData.BaseUrl, cnf.BaseUrl())
		assert.Equal(t, testData.FileStoragePath, cnf.FileStoragePath())
		assert.Equal(t, testData.DatabaseDsn, cnf.DatabaseDsn())
		assert.Equal(t, true, cnf.EnableHTTPS())

		cnf.SetEnableHttps("")

		assert.Equal(t, false, cnf.EnableHTTPS())

	})
}
