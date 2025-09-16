// Package main содержит точку входа для запуска агента сбора метрик.
// Агент собирает метрики и отправляет их на сервер с заданной периодичностью.

package main

import (
	"github.com/alexkozopolianski/go-metrics-tpl/internal/config"
	"github.com/alexkozopolianski/go-metrics-tpl/internal/services"
)

// main — точка входа приложения-агента.
// Получает конфигурацию, создает сервис метрик и запускает процесс сбора и отправки метрик.
func main() {
	cfg := config.GetAgentConfig()               // Получение конфигурации агента
	agent := services.NewAgentMetricService(cfg) // Создание сервиса метрик агента
	agent.Run()                                  // Запуск процесса сбора и отправки метрик
}
