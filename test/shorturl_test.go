package test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/sokol2106/go-url-shortener/internal/cerrors"
	"github.com/sokol2106/go-url-shortener/internal/model"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
	"time"
)

func TestServiceShortURL(t *testing.T) {
	var (
		err    error
		buffer bytes.Buffer
	)

	baseURL := "http://localhost/8080"
	str := storage.NewMemory()
	srvShort := service.NewShortURL(baseURL, str)
	defer srvShort.Close()
	mdl := model.ShortData{OriginalURL: "https://www.youtube.com/", UserID: "111111"}
	mdl1 := model.ShortData{OriginalURL: "https://ya.ru/", UserID: "111111"}
	mdl2 := model.ShortData{OriginalURL: "https://translate.yandex.ru/", UserID: "222222"}
	mdl3 := model.ShortData{OriginalURL: "https://yandex.ru/maps/", UserID: "222222"}

	res1 := make([]service.ResponseUserShortenedURL, 0)
	res2 := make([]service.ResponseUserShortenedURL, 0)

	t.Run("testServiceShortURL", func(t *testing.T) {
		t.Parallel()
		mdl.ShortURL, err = srvShort.AddOriginalURL(mdl.OriginalURL, mdl.UserID)
		require.NoError(t, err)
		res1 = append(res1, service.ResponseUserShortenedURL{OriginalURL: mdl.OriginalURL, ShortURL: mdl.ShortURL})

		mdl1.ShortURL, err = srvShort.AddOriginalURL(mdl1.OriginalURL, mdl1.UserID)
		require.NoError(t, err)
		res1 = append(res1, service.ResponseUserShortenedURL{OriginalURL: mdl1.OriginalURL, ShortURL: mdl1.ShortURL})

		mdl2.ShortURL, err = srvShort.AddOriginalURL(mdl2.OriginalURL, mdl2.UserID)
		require.NoError(t, err)
		res2 = append(res2, service.ResponseUserShortenedURL{OriginalURL: mdl2.OriginalURL, ShortURL: mdl2.ShortURL})

		mdl3.ShortURL, err = srvShort.AddOriginalURL(mdl3.OriginalURL, mdl3.UserID)
		require.NoError(t, err)
		res2 = append(res2, service.ResponseUserShortenedURL{OriginalURL: mdl3.OriginalURL, ShortURL: mdl3.ShortURL})

		err = json.NewEncoder(&buffer).Encode(res1)
		require.NoError(t, err)

		body, err := io.ReadAll(&buffer)
		require.NoError(t, err)

		res, err := srvShort.GetUserShortenedURLs(context.Background(), mdl.UserID)
		require.NoError(t, err)
		assert.JSONEq(t, string(body), string(res))

		err = json.NewEncoder(&buffer).Encode(res2)
		require.NoError(t, err)

		body, err = io.ReadAll(&buffer)
		require.NoError(t, err)

		res, err = srvShort.GetUserShortenedURLs(context.Background(), mdl3.UserID)
		require.NoError(t, err)
		assert.JSONEq(t, string(body), string(res))

	})
}

func TestDeleteURLs(t *testing.T) {
	var err error
	baseURL := ""
	str := storage.NewMemory()
	srvShort := service.NewShortURL(baseURL, str)
	defer srvShort.Close()

	mdl := model.ShortData{OriginalURL: "https://www.youtube.com/", UserID: "111111"}
	mdl1 := model.ShortData{OriginalURL: "https://ya.ru/", UserID: "111111"}
	mdl2 := model.ShortData{OriginalURL: "https://translate.yandex.ru/", UserID: "222222"}
	mdl3 := model.ShortData{OriginalURL: "https://yandex.ru/maps/", UserID: "222222"}

	deleteShortURLs := make([]string, 3)

	t.Run("testDeleteURLs", func(t *testing.T) {
		t.Parallel()
		mdl.ShortURL, err = srvShort.AddOriginalURL(mdl.OriginalURL, mdl.UserID)
		require.NoError(t, err)
		deleteShortURLs[0] = mdl.ShortURL[1:]

		mdl1.ShortURL, err = srvShort.AddOriginalURL(mdl1.OriginalURL, mdl1.UserID)
		require.NoError(t, err)
		deleteShortURLs[1] = mdl1.ShortURL[1:]

		mdl2.ShortURL, err = srvShort.AddOriginalURL(mdl2.OriginalURL, mdl2.UserID)
		require.NoError(t, err)
		deleteShortURLs[2] = mdl2.ShortURL[1:]

		mdl3.ShortURL, err = srvShort.AddOriginalURL(mdl3.OriginalURL, mdl3.UserID)
		require.NoError(t, err)

		srvShort.DeleteOriginalURLs(context.Background(), mdl.UserID, deleteShortURLs)
		time.Sleep(1 * time.Second)

		originalURL, err2 := srvShort.GetOriginalURL(context.Background(), deleteShortURLs[0])
		assert.Equal(t, err2, cerrors.ErrGetShortURLDelete)
		assert.Equal(t, originalURL, "")

		originalURL, err2 = srvShort.GetOriginalURL(context.Background(), deleteShortURLs[1])
		assert.Equal(t, err2, cerrors.ErrGetShortURLDelete)
		assert.Equal(t, originalURL, "")

		originalURL, err2 = srvShort.GetOriginalURL(context.Background(), deleteShortURLs[2])
		require.NoError(t, err2)
		assert.Equal(t, originalURL, mdl2.OriginalURL)

	})
}

func BenchmarkServiceShortURL(b *testing.B) {
	baseURL := "http://localhost/8080"
	str := storage.NewMemory()
	srvShort := service.NewShortURL(baseURL, str)
	defer srvShort.Close()

	for i := 0; i < b.N; i++ {
		sh, _ := srvShort.AddOriginalURL(storage.RandText(20), storage.RandText(3))
		srvShort.GetOriginalURL(context.Background(), sh)
	}

	str.Close()
}
