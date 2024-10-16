// Package server предоставляет возможности для инициализации, запуска и остановки HTTP-сервера.
package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"
)

// Server представляет структуру HTTP-сервера.
type Server struct {
	httpServer *http.Server
}

// NewServer создаёт сервер и возвращает аддрес объекта.
func NewServer(handler http.Handler, addr string) *Server {
	return &Server{
		httpServer: &http.Server{
			Handler: handler,
			Addr:    addr,
		},
	}
}

// Start запускает HTTP-сервер. Сервер начинает слушать входящие запросы.
func (s *Server) Start(enableHTTPS string) error {
	if enableHTTPS != "" {
		return s.httpServer.ListenAndServe()
	}

	return s.httpServer.ListenAndServeTLS(s.CreateCRT())
}

// Stop останавливает HTTP-сервер с возможностью плавного завершения работы.
// ctx - контекст, который может использоваться для задания тайм-аута на завершение работы сервера.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// CreateCRT создаёт сертификат и ключ для HTTPS сервера.
func (s *Server) CreateCRT() (string, string) {

	// cert сертификат сервера для https
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Go Web Server"},
			Country:      []string{"RU"},
		},
		NotAfter:     time.Now(),
		NotBefore:    time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	certPath := "./server.crt"
	keyPath := "./server.key"

	// Сохраняем сертификат в файл
	err = os.WriteFile(certPath, certPEM.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Сохраняем ключ в файл
	err = os.WriteFile(keyPath, privateKeyPEM.Bytes(), 0600) // Даем права только на чтение для владельца
	if err != nil {
		log.Fatal(err)
	}

	// Возвращаем пути к файлам сертификата и ключа
	return certPath, keyPath
}
