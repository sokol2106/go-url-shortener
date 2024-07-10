package test

import (
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestStorageFile(t *testing.T) {
	var (
		req []shorturl.RequestBatch
	)

	req = append(req, shorturl.RequestBatch{CorrelationID: "0", OriginalURL: "testOriginalURL0"})
	req = append(req, shorturl.RequestBatch{CorrelationID: "1", OriginalURL: "testOriginalURL1"})
	original := "testOriginalURL"

	fileName := "testStorage.json"
	objectStorage := storage.NewFile(fileName)

	t.Run("testStorageFile", func(t *testing.T) {
		// Проверка заполнения файла
		short := objectStorage.AddURL(original)
		resOriginal := objectStorage.GetURL(short)
		assert.Equal(t, original, resOriginal)

		resp := objectStorage.AddBatch(req, "")
		original0 := objectStorage.GetURL(strings.ReplaceAll(resp[0].ShortURL, "/", ""))
		original1 := objectStorage.GetURL(strings.ReplaceAll(resp[1].ShortURL, "/", ""))

		assert.Equal(t, req[0].OriginalURL, original0)
		assert.Equal(t, req[1].OriginalURL, original1)

		err := objectStorage.Close()
		require.NoError(t, err)

		// Проверка загрузки из файла
		objectStorage = storage.NewFile(fileName)

		original0 = objectStorage.GetURL(strings.ReplaceAll(resp[0].ShortURL, "/", ""))
		original1 = objectStorage.GetURL(strings.ReplaceAll(resp[1].ShortURL, "/", ""))

		assert.Equal(t, req[0].OriginalURL, original0)
		assert.Equal(t, req[1].OriginalURL, original1)

		err = objectStorage.Close()
		require.NoError(t, err)

		err = os.Remove(fileName)
		require.NoError(t, err)

	})
}

func TestStorageMemory(t *testing.T) {
	var (
		req []shorturl.RequestBatch
	)

	req = append(req, shorturl.RequestBatch{CorrelationID: "0", OriginalURL: "testOriginalURL0"})
	req = append(req, shorturl.RequestBatch{CorrelationID: "1", OriginalURL: "testOriginalURL1"})

	objectStorage := storage.NewMemory()
	original := "testOriginalURL"

	t.Run("testStorageMemory", func(t *testing.T) {
		short := objectStorage.AddURL(original)
		resOriginal := objectStorage.GetURL(short)
		assert.Equal(t, original, resOriginal)

		resp := objectStorage.AddBatch(req, "")
		original0 := objectStorage.GetURL(strings.ReplaceAll(resp[0].ShortURL, "/", ""))
		original1 := objectStorage.GetURL(strings.ReplaceAll(resp[1].ShortURL, "/", ""))

		assert.Equal(t, req[0].OriginalURL, original0)
		assert.Equal(t, req[1].OriginalURL, original1)

		err := objectStorage.Close()
		require.NoError(t, err)

	})
}

func TestStoragePostgresql(t *testing.T) {
	/*	var strg storage.ShortDataList
		strg.Init("")
		db := postgresql.New("host=localhost port=5432 user=postgres password=12345678 dbname=videos sslmode=disable")
		err := db.Connect()
		if err != nil {
			t.Error(err)
			return
		}

		//sh := shorturl.New("http://localhost:8080", strg, db)

		//request := httptest.NewRequest("GET", "/", strings.NewReader("https://practicum.yandex.ru/"))
		//response := httptest.NewRecorder()
		//sh.GetPingDB(response, request)

		err = db.PingContext()
		if err != nil {
			t.Error(err)
			return
		}

		err = db.Disconnect()
		if err != nil {
			t.Error(err)
		}
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
