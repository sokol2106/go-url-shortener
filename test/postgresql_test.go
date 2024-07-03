package test

import (
	"github.com/sokol2106/go-url-shortener/internal/database/postgresql"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"testing"
)

func TestPostgresql(t *testing.T) {
	var strg storage.ShortDataList
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

}
