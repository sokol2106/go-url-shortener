package test

import (
	"context"
	"github.com/sokol2106/go-url-shortener/internal/cerrors"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
)

func TestStorageFile(t *testing.T) {
	var (
		req []service.RequestBatch
	)

	req = append(req, service.RequestBatch{CorrelationID: "0", OriginalURL: "testOriginalURL0"})
	req = append(req, service.RequestBatch{CorrelationID: "1", OriginalURL: "testOriginalURL1"})
	original := "testOriginalURL"

	fileName := "testStorage.json"
	objectStorage := storage.NewFile(fileName)

	t.Run("testStorageFile", func(t *testing.T) {
		// Проверка добавление одной ссылки
		ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Second)
		defer cancel()
		short, err := objectStorage.AddOriginalURL(original, "12345")
		assert.NoError(t, err)
		resOriginal, err := objectStorage.GetOriginalURL(ctx, short)
		require.NoError(t, err)
		assert.Equal(t, original, resOriginal.OriginalURL)

		// Проверка повторного добавление ссылки
		short2, err := objectStorage.AddOriginalURL(original, "12345")
		assert.ErrorIs(t, cerrors.ErrNewShortURL, err)
		assert.Equal(t, short, short2)

		// Проверка добавления массива ссылок
		resp, err := objectStorage.AddOriginalURLBatch(req, "", "12345")
		assert.NoError(t, err)
		original0, err := objectStorage.GetOriginalURL(ctx, strings.ReplaceAll(resp[0].ShortURL, "/", ""))
		require.NoError(t, err)
		original1, err := objectStorage.GetOriginalURL(ctx, strings.ReplaceAll(resp[1].ShortURL, "/", ""))
		require.NoError(t, err)
		assert.Equal(t, req[0].OriginalURL, original0.OriginalURL)
		assert.Equal(t, req[1].OriginalURL, original1.OriginalURL)

		// Проверка повторного добавления массива ссылок
		resp2, err := objectStorage.AddOriginalURLBatch(req, "", "12345")
		assert.ErrorIs(t, cerrors.ErrNewShortURL, err)
		assert.Equal(t, resp, resp2)

		err = objectStorage.Close()
		require.NoError(t, err)

		// Проверка загрузки из файла
		objectStorage = storage.NewFile(fileName)
		original0, err = objectStorage.GetOriginalURL(ctx, strings.ReplaceAll(resp[0].ShortURL, "/", ""))
		require.NoError(t, err)
		original1, err = objectStorage.GetOriginalURL(ctx, strings.ReplaceAll(resp[1].ShortURL, "/", ""))
		require.NoError(t, err)
		assert.Equal(t, req[0].OriginalURL, original0.OriginalURL)
		assert.Equal(t, req[1].OriginalURL, original1.OriginalURL)

		err = objectStorage.Close()
		require.NoError(t, err)
		err = os.Remove(fileName)
		require.NoError(t, err)

	})
}

func TestStorageMemory(t *testing.T) {
	var (
		req []service.RequestBatch
	)

	req = append(req, service.RequestBatch{CorrelationID: "0", OriginalURL: "testOriginalURL0"})
	req = append(req, service.RequestBatch{CorrelationID: "1", OriginalURL: "testOriginalURL1"})

	objectStorage := storage.NewMemory()
	original := "testOriginalURL"

	t.Run("testStorageMemory", func(t *testing.T) {
		// Проверка добавление одной ссылки
		ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Second)
		defer cancel()
		short, err := objectStorage.AddOriginalURL(original, "123456")
		assert.NoError(t, err)
		resOriginal, err := objectStorage.GetOriginalURL(ctx, short)
		require.NoError(t, err)
		assert.Equal(t, original, resOriginal.OriginalURL)

		// Проверка повторного добавление ссылки
		short2, err := objectStorage.AddOriginalURL(original, "123456")
		assert.ErrorIs(t, cerrors.ErrNewShortURL, err)
		assert.Equal(t, short, short2)

		// Проверка добавления массива ссылок
		resp, err := objectStorage.AddOriginalURLBatch(req, "", "123456")
		assert.NoError(t, err)
		original0, err := objectStorage.GetOriginalURL(ctx, strings.ReplaceAll(resp[0].ShortURL, "/", ""))
		require.NoError(t, err)
		original1, err := objectStorage.GetOriginalURL(ctx, strings.ReplaceAll(resp[1].ShortURL, "/", ""))
		require.NoError(t, err)
		assert.Equal(t, req[0].OriginalURL, original0.OriginalURL)
		assert.Equal(t, req[1].OriginalURL, original1.OriginalURL)

		// Проверка повторного добавления массива ссылок
		resp2, err := objectStorage.AddOriginalURLBatch(req, "", "123456")
		assert.ErrorIs(t, cerrors.ErrNewShortURL, err)
		assert.Equal(t, resp, resp2)

		err = objectStorage.Close()
		require.NoError(t, err)

	})
}

func TestStoragePostgresql(t *testing.T) {
	/*store := storage.NewPostgresql("host=localhost port=5432 user=postgres password=12345678 dbname=test sslmode=disable")
	err := store.Connect()
	defer store.Close()
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("testStoragePostgresql", func(t *testing.T) {
		sh := shorturl.New("http://localhost:8080", store)

		request := httptest.NewRequest("Post", "/", strings.NewReader("https://practicum.yandex.ru/"))
		response := httptest.NewRecorder()
		sh.Post(response, request)

		result := response.Result()
		assert.Equal(t, http.StatusCreated, result.StatusCode)

	})

	*/
}

func TestPostgresqlMocks(t *testing.T) {
	/*	ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock_shorturl.NewMockDatabase(ctrl)
		db.EXPECT().PingContext().Return(nil)

		strg := memoryStorage.New()
		sh := shorturl.New("http://localhost:8080", strg, db)
		t.Run("Test ping mocks", func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", strings.NewReader(""))
			response := httptest.NewRecorder()
			sh.GetPing(response, request)

			assert.Equal(t, http.StatusOK, response.Code)

			db.EXPECT().PingContext().Return(errors.New("errr"))
			request = httptest.NewRequest("GET", "/", strings.NewReader(""))
			response = httptest.NewRecorder()
			sh.GetPing(response, request)

			assert.Equal(t, http.StatusInternalServerError, response.Code)
		})

	*/
}
