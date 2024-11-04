package test

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

import (
	"github.com/sokol2106/go-url-shortener/internal/handlers"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type ServerTestSuite struct {
	suite.Suite
	server *httptest.Server
	cookie *http.Cookie
}

func (suite *ServerTestSuite) SetupSuite() {
	objStorage := storage.NewMemory()
	srvShortURL := service.NewShortURL("http://localhost:8080", objStorage)
	handler := handlers.NewHandlers(srvShortURL, "192.168.1.0/24")
	suite.server = httptest.NewServer(handlers.Router(handler))
	srvShortURL.SetRedirectURL(suite.server.URL)
}

func (suite *ServerTestSuite) TearSuiteDownSuite() {
	suite.server.Close()
}

func (suite *ServerTestSuite) TestPostAndGet() {
	resp, err := http.Post(suite.server.URL+"/", "text/plain", strings.NewReader("https://www.postgresql.org/"))
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	require.Len(suite.T(), resp.Cookies(), 1)
	suite.cookie = resp.Cookies()[0]
	resBody, err := io.ReadAll(resp.Body)
	require.NoError(suite.T(), err)
	resp.Body.Close()

	resp, err = http.Get(string(resBody))
	require.NoError(suite.T(), err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *ServerTestSuite) TestCoockie() {
	resp, err := http.Post(suite.server.URL+"/", "text/plain", strings.NewReader("https://github.com/"))
	require.NoError(suite.T(), err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	cookies := resp.Cookies()
	assert.Len(suite.T(), cookies, 1)
	assert.Equal(suite.T(), "user", cookies[0].Name)
}

func (suite *ServerTestSuite) TestDeleteUserShortenedURLs() {
	resp, err := http.Post(suite.server.URL+"/", "text/plain", strings.NewReader("https://yandex.ru/maps/"))
	require.NoError(suite.T(), err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	require.Len(suite.T(), resp.Cookies(), 1)
	suite.cookie = resp.Cookies()[0]

	reqBody := "[\"rfgtyhju\",\"qazxswed\",\"wsxcderf\"]"
	req, err := http.NewRequest(http.MethodDelete, suite.server.URL+"/api/user/urls", strings.NewReader(reqBody))
	require.NoError(suite.T(), err)
	req.AddCookie(suite.cookie)
	client := &http.Client{}
	resp, err = client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusAccepted, resp.StatusCode)

}

func (suite *ServerTestSuite) TestPostJSON() {
	reqBody := "{\"url\": \"https://practicum.yandex.ru\"}"
	resp, err := http.Post(suite.server.URL+"/api/shorten", "application/json", strings.NewReader(reqBody))
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	require.Equal(suite.T(), "application/json", resp.Header.Get("Content-Type"))

	resBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	require.NoError(suite.T(), err)

	var respJS handlers.ResponseJSON
	err = json.Unmarshal(resBody, &respJS)
	require.NoError(suite.T(), err)

	resp, err = http.Get(respJS.Result)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

}

func (suite *ServerTestSuite) TestGzipCompression() {
	encodeBuffer := bytes.NewBuffer(nil)
	writeGZip := gzip.NewWriter(encodeBuffer)

	_, err := writeGZip.Write([]byte("https://about.gitlab.com/"))
	require.NoError(suite.T(), err)

	err = writeGZip.Close()
	require.NoError(suite.T(), err)

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, suite.server.URL, encodeBuffer)
	require.NoError(suite.T(), err)
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(request)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	require.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	require.Equal(suite.T(), "text/plain", resp.Header.Get("Content-Type"))

	readGZip, err := gzip.NewReader(resp.Body)
	require.NoError(suite.T(), err)

	bodyDecode, err := io.ReadAll(readGZip)
	require.NoError(suite.T(), err)

	resp2, err := http.Get(string(bodyDecode))
	require.NoError(suite.T(), err)
	defer resp2.Body.Close()
	assert.Equal(suite.T(), http.StatusOK, resp2.StatusCode)
}

func (suite *ServerTestSuite) TestStats() {
	resp, err := http.Post(suite.server.URL+"/", "text/plain", strings.NewReader("https://www.postgresql3.org/"))
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	resp, err = http.Post(suite.server.URL+"/", "text/plain", strings.NewReader("https://www.postgresql2.org/"))
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	req, err := http.NewRequest("GET", suite.server.URL+"/api/internal/stats", nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("X-Real-IP", "192.168.1.10")
	client := &http.Client{}
	resp, err = client.Do(req)

	require.NoError(suite.T(), err)
	require.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	resBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	var respJS handlers.ResponseStats
	err = json.Unmarshal(resBody, &respJS)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 7, respJS.Urls)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func BenchmarkSuite(b *testing.B) {

	suite := new(ServerTestSuite)
	// Инициализация suite
	suite.SetupSuite() // Используем SetupTest, как в тестах

	for i := 0; i < b.N; i++ {
		resp, _ := http.Post(suite.server.URL+"/", "text/plain", strings.NewReader(storage.RandText(30)))
		suite.cookie = resp.Cookies()[0]
		resBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		resp, _ = http.Get(string(resBody))
		resp.Body.Close()
	}
}

func BenchmarkHundredSuite(b *testing.B) {

	suite := new(ServerTestSuite)
	// Инициализация suite
	suite.SetupSuite() // Используем SetupTest, как в тестах

	for i := 0; i < b.N; i++ {
		resp, _ := http.Post(suite.server.URL+"/", "text/plain", strings.NewReader(storage.RandText(100)))
		suite.cookie = resp.Cookies()[0]
		resBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		resp, _ = http.Get(string(resBody))
		resp.Body.Close()
	}
}

func BenchmarkThousandSuite(b *testing.B) {

	suite := new(ServerTestSuite)
	// Инициализация suite
	suite.SetupSuite() // Используем SetupTest, как в тестах

	for i := 0; i < b.N; i++ {
		resp, _ := http.Post(suite.server.URL+"/", "text/plain", strings.NewReader(storage.RandText(1000)))
		suite.cookie = resp.Cookies()[0]
		resBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		resp, _ = http.Get(string(resBody))
		resp.Body.Close()
	}
}
