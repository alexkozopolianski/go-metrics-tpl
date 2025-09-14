package handler

import "github.com/alexkozopolianski/go-metrics-tpl/internal/domain"

type Storager interface {
	Save(metricType string, metricName string, value any) error
	Get(metricType, metricName string) (domain.Metric, bool)
	GetAll() []domain.Metric
}
