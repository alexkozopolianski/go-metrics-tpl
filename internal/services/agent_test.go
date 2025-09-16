package services

import (
	"testing"

	"github.com/Elvilius/go-musthave-metrics-tpl/internal/config"
	"github.com/stretchr/testify/assert"
)

func Test_GetMetrics(t *testing.T) {
	testCfg := config.AgentConfig{PollInterval: 10, ServerAddress: "localhost:8080", ReportInterval: 3}
	agentServiceMetrics := NewAgentMetricService(testCfg)

	metrics := agentServiceMetrics.GetMetric()
	expectedGauges := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
		"StackSys", "Sys", "TotalAlloc", "RandomValue",
	}

	expectedCounters := []string{"PollCount"}

	t.Run("expected Gauges", func(t *testing.T) {
		for _, expectedGauge := range expectedGauges {
			_, ok := metrics[expectedGauge]
			assert.True(t, ok)
		}
	})

	t.Run("expected Counters", func(t *testing.T) {
		for _, expectedCounter := range expectedCounters {
			_, ok := metrics[expectedCounter]
			assert.True(t, ok)
		}
	})
}
