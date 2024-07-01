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
	var shortDataList storage.ShortDataList

	fileName := "testStorage.json"
	shortDataList.Init(fileName)
	shortDataTest := model.ShortData{UUID: "testUUID", OriginalURL: "testOriginalURL"}

	t.Run("testStorage", func(t *testing.T) {
		// Проверяем сокращение и получение URL
		shortDataTest.ShortURL = shortDataList.AddURL(shortDataTest.OriginalURL)
		assert.Equal(t, shortDataTest.OriginalURL, shortDataList.GetURL(shortDataTest.ShortURL))

	})

	err := shortDataList.Close()
	require.NoError(t, err)

	err = os.Remove(fileName)
	require.NoError(t, err)

}
