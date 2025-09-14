package server

import (
	"net/http"

	"github.com/alexkozopolianski/go-metrics-tpl/internal/config"
	handler "github.com/alexkozopolianski/go-metrics-tpl/internal/handlers"
	middleware "github.com/alexkozopolianski/go-metrics-tpl/internal/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	handler *handler.Handler
	router  *chi.Mux
	cfg     *config.ServerConfig
	logger  *zap.SugaredLogger
}

func New(cfg *config.ServerConfig, handler *handler.Handler, logger *zap.SugaredLogger) *Server {
	router := chi.NewRouter()

	server := &Server{handler: handler, router: router, cfg: cfg, logger: logger}

	router.Get("/", server.handler.All)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", server.handler.Update)
	router.Get("/value/{metricType}/{metricName}", server.handler.Value)

	return server
}

func (s *Server) Run() {
	err := http.ListenAndServe(s.cfg.Address, middleware.Logging(s.router, s.logger))
	if err != nil {
		panic(err)
	}
}
