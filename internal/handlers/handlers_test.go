package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexkozopolianski/go-metrics-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Update(t *testing.T) {
	type want struct {
		status int
	}

	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name:    "positive test #1",
			request: "/update/gauge/cpu/7513",
			want: want{
				status: 200,
			},
		},
		{
			name:    "positive test #2",
			request: "/update/counter/cpu/8",
			want: want{
				status: 200,
			},
		},
		{
			name:    "negative test #1",
			request: "/update/test123123123/cpu/8",
			want: want{
				status: 400,
			},
		},
		{
			name:    "negative test #2",
			request: "/update/",
			want: want{
				status: 404,
			},
		},
	}
	for _, tt := range tests {
		memStorage := storage.NewMemStorage()
		h := NewHandler(memStorage)
		router := chi.NewRouter()
		router.Post("/update/{metricType}/{metricName}/{metricValue}", h.Update)

		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)
			result.Body.Close()
		})
	}
}

func TestHandler_Value(t *testing.T) {
	type want struct {
		status int
	}

	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name:    "positive test #1",
			request: "/value/gauge/Alloc",
			want: want{
				status: 200,
			},
		},
		{
			name:    "positive test #2",
			request: "/value/counter/PollCount",
			want: want{
				status: 200,
			},
		},
		{
			name:    "negative test #1",
			request: "/value/test/Alloc",
			want: want{
				status: 404,
			},
		},
		{
			name:    "negative test #2",
			request: "/value/counter/test",
			want: want{
				status: 404,
			},
		},
	}
	memStorage := storage.NewMemStorage()
	memStorage.Save("gauge", "Alloc", "1.1")
	memStorage.Save("counter", "PollCount", "100")

	h := NewHandler(memStorage)
	router := chi.NewRouter()
	router.Get("/value/{metricType}/{metricName}", h.Value)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.status, result.StatusCode)
			result.Body.Close()
		})
	}
}
