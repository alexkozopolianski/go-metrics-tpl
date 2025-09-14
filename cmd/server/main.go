package main

import (
	handler "github.com/alexkozopolianski/go-metrics-tpl/internal/handlers"
	"github.com/alexkozopolianski/go-metrics-tpl/internal/storage"
	"go.uber.org/zap"

	"github.com/alexkozopolianski/go-metrics-tpl/internal/config"
	"github.com/alexkozopolianski/go-metrics-tpl/internal/server"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	sugarLogger := logger.Sugar()

	memStorage := storage.NewMemStorage()
	handler := handler.NewHandler(memStorage)
	cfg := config.GetServerConfig()
	server := server.New(&cfg, handler, sugarLogger)

	server.Run()
}
