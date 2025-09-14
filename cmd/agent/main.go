package main

import (
	"github.com/alexkozopolianski/go-metrics-tpl/internal/config"
	"github.com/alexkozopolianski/go-metrics-tpl/internal/services"
)

func main() {
	cfg := config.GetAgentConfig()

	agentServiceMetrics := services.NewAgentMetricService(cfg)
	agentServiceMetrics.SendMetrics()
}
