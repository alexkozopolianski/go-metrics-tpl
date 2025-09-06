package storage

import (
	"strconv"

	"github.com/alexkozopolianski/go-metrics-tpl/internal/domain"
)

type MemStorage struct {
	metrics map[string]domain.Metric
}

func NewMemStorage() Storage {
	return &MemStorage{metrics: make(map[string]domain.Metric)}
}

func (r *MemStorage) Save(metricType string, metricName string, value any) error {
	existMeric, ok := r.Get(metricType, metricName)

	if metricType == domain.Gauge {
		parsedValueFloat, err := strconv.ParseFloat(value.(string), 64)

		if err != nil {
			return err
		}
		r.metrics[metricName] = domain.Metric{Type: metricType, Name: metricName, Value: parsedValueFloat}
		return nil
	}
	if metricType == domain.Counter {
		parsedValue, err := strconv.ParseInt(value.(string), 10, 64)

		if err != nil {
			return err
		}
		if !ok {
			r.metrics[metricName] = domain.Metric{Type: metricType, Name: metricName, Value: parsedValue}
			return nil

		} else {
			existMeric.Value = existMeric.Value.(int64) + parsedValue
			r.metrics[metricName] = existMeric
		}
	}
	return nil
}

func (r *MemStorage) Get(metricType, metricName string) (domain.Metric, bool) {
	metric, ok := r.metrics[metricName]

	if !ok {
		return domain.Metric{}, false
	}
	if metric.Type != metricType {
		return domain.Metric{}, false
	}
	return metric, true
}

func (r *MemStorage) GetAll() []domain.Metric {
	all := make([]domain.Metric, 0, len(r.metrics))

	for _, metric := range r.metrics {
		all = append(all, metric)
	}
	return all
}
