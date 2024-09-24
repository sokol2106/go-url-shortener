// Package middleware предоставляет middleware для обработки HTTP-запросов и ответов,
// включая логирование через zap.NewDevelopment.
package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (

	// responseData представляет структуру для хранения информации о ответе.
	// status - код статуса HTTP-ответа.
	// size - размер HTTP-ответа в байтах.
	responseData struct {
		status int
		size   int
	}

	// loggingResponseWriter представляет обёртку над http.ResponseWriter,
	// которая захватывает данные о статусе и размере ответа.
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write записывает данные в HTTP-ответ и обновляет размер ответа.
// Возвращает количество записанных байт и ошибку, если она произошла.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

// WriteHeader отправляет код статуса в ответ и обновляет статус ответа.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

// LoggingResponseRequest является middleware, которая логирует HTTP-запросы и ответы.
// Принимает http.Handler и возвращает новый http.Handler с логированием.
// Логируется время выполения запроса duration.
func LoggingResponseRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		//response
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		handler.ServeHTTP(&lw, r)

		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}

		defer logger.Sync()
		sugar := logger.Sugar()

		duration := time.Since(start)
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
		)

		sugar.Infoln(
			"status", responseData.status,
			"size", responseData.size,
		)
	})

}
