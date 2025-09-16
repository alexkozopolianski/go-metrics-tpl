// Package main содержит точку входа для запуска HTTP-сервера метрик.
// Здесь инициализируются логгер, хранилище, обработчики, конфиг и сервер.

package main

import (
	handler "github.com/alexkozopolianski/go-metrics-tpl/internal/handlers"
	"github.com/alexkozopolianski/go-metrics-tpl/internal/storage"
	"go.uber.org/zap"

	"github.com/alexkozopolianski/go-metrics-tpl/internal/config"
	"github.com/alexkozopolianski/go-metrics-tpl/internal/server"
)

// main — точка входа приложения. Инициализирует все зависимости и запускает сервер.
func main() {
	logger, err := zap.NewDevelopment() // Инициализация логгера
	if err != nil {
		panic(err)
	}

	defer func() {
		err := logger.Sync() // Синхронизация логгера перед завершением
		if err != nil {
			panic(err)
		}
	}()

	sugarLogger := logger.Sugar() // Упрощённый логгер

	memStorage := storage.NewMemStorage()            // Создание хранилища метрик в памяти
	handler := handler.NewHandler(memStorage)        // Создание обработчиков с хранилищем
	cfg := config.GetServerConfig()                  // Получение конфигурации сервера
	server := server.New(&cfg, handler, sugarLogger) // Создание сервера

	server.Run() // Запуск HTTP-сервера
}
