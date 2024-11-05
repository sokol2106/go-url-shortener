package grpchandlers

import (
	"context"
	"encoding/json"
	"github.com/sokol2106/go-url-shortener/internal/handlers"
	"github.com/sokol2106/go-url-shortener/internal/handlers/grpchandlers/grpcservice/proto"
	"github.com/sokol2106/go-url-shortener/internal/service"
)

// URLShortenerServer предоставляет gRPC-сервис сокращения URL.
type URLShortenerServer struct {
	proto.UnimplementedURLShortenerServer
	srvShortURL      *service.ShortURL      // Сервис сокращения URL
	srvAuthorization *service.Authorization // Сервис авторизации
	trustedSubnet    string                 // Конфигурация сервера
}

// NewURLShortenerServer создает новый экземпляр gRPC-сервера для сокращения URL.
func NewURLShortenerServer(srvSh *service.ShortURL, srvAu *service.Authorization, subnet string) *URLShortenerServer {
	return &URLShortenerServer{
		srvShortURL:      srvSh,
		srvAuthorization: srvAu,
		trustedSubnet:    subnet,
	}
}

// CreateShortURL обрабатывает запрос на создание сокращенного URL.
func (s *URLShortenerServer) CreateShortURL(ctx context.Context, req *proto.CreateShortURLRequest) (*proto.CreateShortURLResponse, error) {
	id, err := s.srvAuthorization.GetUserID(req.Token)
	if err != nil {
		return nil, err
	}

	shortURL, err := s.srvShortURL.AddOriginalURL(req.Url, id)
	if err != nil {
		return nil, err
	}

	return &proto.CreateShortURLResponse{Result: shortURL}, nil
}

// CreateShortURL обрабатывает JSON запрос на создание сокращенного URL.
func (s *URLShortenerServer) CreateShortURLJSON(ctx context.Context, req *proto.CreateShortURLRequestJSON) (*proto.CreateShortURLResponseJSON, error) {
	var (
		reqJS handlers.RequestJSON
		resJS handlers.ResponseJSON
	)

	id, err := s.srvAuthorization.GetUserID(req.Token)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(req.Url), &reqJS)

	if err != nil {
		return nil, err
	}

	resJS.Result, err = s.srvShortURL.AddOriginalURL(req.Url, id)
	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(resJS)
	if err != nil {
		return nil, err
	}

	return &proto.CreateShortURLResponseJSON{Result: string(response)}, nil
}

// GetOriginalURL обрабатывает запрос на получение оригинального URL по сокращенному.
func (s *URLShortenerServer) GetOriginalURL(ctx context.Context, req *proto.GetOriginalURLRequest) (*proto.GetOriginalURLResponse, error) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	originalURL, err := s.srvShortURL.GetOriginalURL(ctx2, req.Path)
	if err != nil {
		return nil, err
	}
	return &proto.GetOriginalURLResponse{OriginalUrl: originalURL}, nil
}

// GetUserShortenedURLs обрабатывает запрос на получение всех сокращенных URL пользователя.
func (s *URLShortenerServer) GetUserShortenedURLs(ctx context.Context, req *proto.GetUserShortenedURLsRequest) (*proto.GetUserShortenedURLsResponse, error) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	id, err := s.srvAuthorization.GetUserID(req.Token)
	if err != nil {
		return nil, err
	}

	URLs, err := s.srvShortURL.GetUserShortenedURLs(ctx2, id)
	if err != nil {
		return nil, err
	}

	strURLs := make([]string, len(URLs))

	for i, b := range URLs {
		strURLs[i] = string(b)
	}

	return &proto.GetUserShortenedURLsResponse{Urls: strURLs}, nil
}

// DeleteUserShortenedURLs обрабатывает запрос на удаление сокращенных URL пользователя.
func (s *URLShortenerServer) DeleteUserShortenedURLs(ctx context.Context, req *proto.DeleteUserShortenedURLsRequest) (*proto.DeleteUserShortenedURLsResponse, error) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	id, err := s.srvAuthorization.GetUserID(req.Token)
	if err != nil {
		return nil, err
	}

	s.srvShortURL.DeleteOriginalURLs(ctx2, id, req.Urls)
	return &proto.DeleteUserShortenedURLsResponse{}, nil
}
