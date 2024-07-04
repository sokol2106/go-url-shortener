package test

import (
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig(t *testing.T) {
	testData := struct {
		URL  string
		Addr string
		Port string
		Host string
	}{
		URL:  "http://127.0.0.1:8080",
		Addr: "127.0.0.1:8080",
		Port: "8080",
		Host: "127.0.0.1",
	}

	t.Run("testConfig", func(t *testing.T) {
		cnf, err := config.NewConfigURL(testData.URL)
		require.NoError(t, err)

		assert.Equal(t, testData.Port, cnf.Port())
		assert.Equal(t, testData.Addr, cnf.Addr())
		assert.Equal(t, testData.Host, cnf.Host())
		assert.Equal(t, testData.URL, cnf.URL())
	})
}
