// Package server реализует HTTP-сервер для обработки запросов к метрикам.
// Сервер настраивает роутер, регистрирует обработчики и запускает HTTP-сервис.

package server

import (
	"net/http"

	"github.com/alexkozopolianski/go-metrics-tpl/internal/config"
	handler "github.com/alexkozopolianski/go-metrics-tpl/internal/handlers"
	"github.com/alexkozopolianski/go-metrics-tpl/internal/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Server — структура, инкапсулирующая обработчики, роутер, конфиг и логгер.
type Server struct {
	handler *handler.Handler     // Обработчик HTTP-запросов
	router  *chi.Mux             // HTTP-роутер
	cfg     *config.ServerConfig // Конфигурация сервера
	logger  *zap.SugaredLogger   // Логгер
}

// New создает и настраивает новый экземпляр Server с роутером и обработчиками.
// Регистрирует все необходимые маршруты для работы с метриками.
func New(cfg *config.ServerConfig, handler *handler.Handler, logger *zap.SugaredLogger) *Server {
	router := chi.NewRouter()

	server := &Server{handler: handler, router: router, cfg: cfg, logger: logger}

	// Регистрируем маршруты для работы с метриками
	router.Get("/", server.handler.All)                               // Получить все метрики
	router.Post("/update/{type}/{id}/{value}", server.handler.Update) // Обновить метрику через URL
	router.Post("/update/", server.handler.UpdateJSON)                // Обновить метрику через JSON
	router.Post("/value", server.handler.ValueJSON)                   // Получить метрику через JSON
	router.Post("/value/", server.handler.ValueJSON)                  // Получить метрику через JSON (альтернативный путь)
	router.Get("/value/{type}/{id}", server.handler.Value)            // Получить метрику по типу и id

	return server
}

// Run запускает HTTP-сервер на указанном в конфиге адресе.
// Использует middleware для логирования запросов.
func (s *Server) Run() {
	err := http.ListenAndServe(s.cfg.Address, middleware.Logging(s.router, s.logger))
	if err != nil {
		panic(err)
	}
}
