// Package middleware содержит middleware для HTTP-сервера,
// в частности — middleware для логирования запросов.

package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// responseData хранит информацию о статусе и размере ответа.
type (
	responseData struct {
		status int // HTTP-статус ответа
		size   int // Размер тела ответа в байтах
	}

	// loggingResponseWriter оборачивает http.ResponseWriter для сбора информации о статусе и размере ответа.
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Logging — middleware для логирования HTTP-запросов и ответов.
// Логирует URI, метод, статус, длительность обработки и размер ответа.
func Logging(h http.Handler, logger *zap.SugaredLogger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{size: 0, status: 0}

		lw := &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(lw, r)

		duration := time.Since(start)

		logger.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", lw.responseData.status,
			"duration", duration,
			"size", lw.responseData.size,
		)
	}

	return http.HandlerFunc(fn)
}

// Write реализует интерфейс http.ResponseWriter и считает размер записанных данных.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader реализует интерфейс http.ResponseWriter и сохраняет статус ответа.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.responseData.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
