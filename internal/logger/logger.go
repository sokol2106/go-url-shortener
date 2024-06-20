package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

var sugar zap.SugaredLogger

func Init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func LoggingResponseRequest(handler http.Handler) http.HandlerFunc {
	logFn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		sugar = *logger.Sugar()

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

	return logFn
}
