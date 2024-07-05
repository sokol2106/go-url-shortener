package test

import (
	"github.com/sokol2106/go-url-shortener/internal/model"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestStorage(t *testing.T) {
	fileName := "testStorage.json"
	objectStorage := storage.NewFile(fileName)
	shortDataTest := model.ShortData{UUID: "testUUID", OriginalURL: "testOriginalURL"}

	t.Run("testStorage", func(t *testing.T) {
		// Проверяем сокращение и получение URL
		shortDataTest.ShortURL = objectStorage.AddURL(shortDataTest.OriginalURL)
		assert.Equal(t, shortDataTest.OriginalURL, objectStorage.GetURL(shortDataTest.ShortURL))

	})

	err := objectStorage.Close()
	require.NoError(t, err)

	err = os.Remove(fileName)
	require.NoError(t, err)
}
