// Package config содержит функции и структуры для конфигурирования агента и сервера.
// Позволяет получать параметры из переменных окружения и флагов командной строки.

package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// AgentConfig содержит параметры конфигурации для агента.
type AgentConfig struct {
	ServerAddress  string // Адрес сервера для отправки метрик
	PollInterval   int    // Интервал сбора метрик (сек)
	ReportInterval int    // Интервал отправки метрик на сервер (сек)
}

// ServerConfig содержит параметры конфигурации для сервера.
type ServerConfig struct {
	Address string // Адрес, на котором запускается сервер
}

// getEnvOrDefaultString возвращает значение переменной окружения envVar,
// либо defaultValue, если переменная не установлена.
func getEnvOrDefaultString(envVar string, defaultValue string) string {
	if value, ok := os.LookupEnv(envVar); ok {
		return value
	}
	return defaultValue
}

// getEnvOrDefaultInt возвращает значение переменной окружения envVar как int,
// либо defaultValue, если переменная не установлена или не может быть преобразована.
func getEnvOrDefaultInt(envVar string, defaultValue int) int {
	if value, ok := os.LookupEnv(envVar); ok {
		if parsedValue, err := strconv.Atoi(value); err == nil {
			return parsedValue
		}
	}
	return defaultValue
}

// GetAgentConfig возвращает конфигурацию агента.
// Приоритет: переменные окружения → флаги командной строки → значения по умолчанию.
func GetAgentConfig() AgentConfig {
	cfg := AgentConfig{
		ServerAddress:  getEnvOrDefaultString("ADDRESS", "localhost:8080"),
		PollInterval:   getEnvOrDefaultInt("POLL_INTERVAL", 3),
		ReportInterval: getEnvOrDefaultInt("REPORT_INTERVAL", 10),
	}

	pollInterval := flag.Int("p", cfg.PollInterval, "pollInterval")
	reportInterval := flag.Int("r", cfg.ReportInterval, "reportInterval")
	serverAddress := flag.String("a", cfg.ServerAddress, "server address")
	flag.Parse()

	cfg.PollInterval = *pollInterval
	cfg.ReportInterval = *reportInterval
	cfg.ServerAddress = *serverAddress

	fmt.Println("Server Address:", cfg.ServerAddress)
	fmt.Println("Report Interval:", cfg.ReportInterval)
	fmt.Println("Poll Interval:", cfg.PollInterval)

	return cfg
}

// GetServerConfig возвращает конфигурацию сервера.
// Приоритет: переменная окружения ADDRESS → флаг командной строки → значение по умолчанию.
func GetServerConfig() ServerConfig {
	cfg := ServerConfig{
		Address: getEnvOrDefaultString("ADDRESS", "localhost:8080"),
	}
	serverAddress := flag.String("a", cfg.Address, "server address")
	flag.Parse()

	cfg.Address = *serverAddress

	fmt.Println("Server Address:", cfg.Address)
	return cfg
}
