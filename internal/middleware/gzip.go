// Package middleware предоставляет middleware для обработки HTTP-запросов и ответов,
// включая сжатие содержимого через gzip.
package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

// compressResponseWriter представляет обёртку над http.ResponseWriter, которая сжимает
// исходящие данные с использованием gzip.
type compressResponseWriter struct {
	rw http.ResponseWriter
	zp *gzip.Writer
}

// compressReader представляет обёртку над io.ReadCloser, которая распаковывает входящие
// сжатые данные с использованием gzip.
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// newCompressResponseWriter создает новый compressResponseWriter, который будет сжимать ответы с использованием gzip.
// Принимает оригинальный http.ResponseWriter и возвращает ссылку на compressResponseWriter.
func newCompressResponseWriter(w http.ResponseWriter) *compressResponseWriter {
	return &compressResponseWriter{
		rw: w,
		zp: gzip.NewWriter(w),
	}
}

// Header возвращает заголовки HTTP-ответа из оригинального ResponseWriter.
func (c *compressResponseWriter) Header() http.Header {
	return c.rw.Header()
}

// Write записывает сжатые данные в поток ответа.
// Возвращает количество записанных байт и ошибку, если она произошла.
func (c *compressResponseWriter) Write(p []byte) (int, error) {
	return c.zp.Write(p)
}

// WriteHeader отправляет HTTP-статус в ответ, добавляя заголовок "Content-Encoding: gzip"
// для успешных ответов (код меньше 300).
func (c *compressResponseWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.rw.Header().Set("Content-Encoding", "gzip")
	}

	c.rw.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer, завершив сжатие данных.
func (c *compressResponseWriter) Close() error {
	return c.zp.Close()
}

// newCompressReader создает новый compressReader, который распаковывает сжатые данные
// с использованием gzip. Принимает io.ReadCloser и возвращает compressReader и ошибку (если есть).
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read считывает распакованные данные из потока.
// Возвращает количество считанных байт и ошибку, если она произошла.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close закрывает как оригинальный ReadCloser, так и gzip.Reader.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// СompressionResponseRequest является middleware, который обрабатывает сжатие/распаковку данных
// для HTTP-запросов и ответов с использованием gzip. Если клиент поддерживает gzip,
// данные ответа сжимаются. Если запрос содержит сжатые данные, они распаковываются перед обработкой.
func СompressionResponseRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := w

		// Write
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressResponseWriter(w)
			response = cw
			defer cw.Close()
		}

		// Read
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				log.Printf("CompressReader error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		handler.ServeHTTP(response, r)
	})
}
