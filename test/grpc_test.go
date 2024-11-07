package test

import (
	"context"
	"github.com/sokol2106/go-url-shortener/internal/handlers/grpchandlers"
	"github.com/sokol2106/go-url-shortener/internal/handlers/grpchandlers/grpcservice/proto"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"strings"
	"testing"
)

var listener *bufconn.Listener

func init() {
	listener = bufconn.Listen(1024 * 1024) // создаем буфер для соединений
}

func TestGRPC(t *testing.T) {
	objStorage := storage.NewMemory()
	srvShortURL := service.NewShortURL("", objStorage)

	// Создаем сервер
	grpcServer := grpc.NewServer()
	urlShortenerServer := grpchandlers.NewURLShortenerServer(srvShortURL, "")
	proto.RegisterURLShortenerServer(grpcServer, urlShortenerServer)

	lis, err := net.Listen("tcp", ":3200")
	if err != nil {
		t.Error(err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewURLShortenerClient(conn)

	t.Run("testGRPC", func(t *testing.T) {
		token, err := srvShortURL.GetAuthorization().NewUserToken()
		require.NoError(t, err)

		req := &proto.CreateShortURLRequest{
			Token: token,
			Url:   "https://yandex.ru",
		}

		res1, err := client.CreateShortURL(context.Background(), req)
		require.NoError(t, err)

		reqJSON := &proto.CreateShortURLRequestJSON{
			Url:   "{\"urltt\": \"https://practicum.yandex.ru\"}",
			Token: token,
		}

		_, err = client.CreateShortURLJSON(context.Background(), reqJSON)
		require.NoError(t, err)

		reqGet := &proto.GetOriginalURLRequest{
			Path: res1.Result[1:],
		}

		resGet, err := client.GetOriginalURL(context.Background(), reqGet)
		require.NoError(t, err)
		assert.Equal(t, resGet.OriginalUrl, "https://yandex.ru")

		reqGetU := &proto.GetUserShortenedURLsRequest{
			Token: token,
		}

		resGetU, err := client.GetUserShortenedURLs(context.Background(), reqGetU)
		require.NoError(t, err)
		result := strings.Join(resGetU.Urls, " ")
		assert.NotEmpty(t, result)

	})

}
