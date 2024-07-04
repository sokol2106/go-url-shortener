package test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	mock_shorturl "github.com/sokol2106/go-url-shortener/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostgresql(t *testing.T) {
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mock_shorturl.NewMockDatabase(ctrl)
	db.EXPECT().PingContext().Return(nil)

	var strg storage.ShortDataList
	strg.Init("")
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
}
