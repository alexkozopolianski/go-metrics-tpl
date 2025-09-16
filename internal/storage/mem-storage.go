package storage

import (
	"fmt"

	handler "github.com/alexkozopolianski/go-metrics-tpl/internal/handlers"
	"github.com/alexkozopolianski/go-metrics-tpl/internal/models"
)

// MemStorage реализует интерфейс хранилища метрик в оперативной памяти.
// Метрики сохраняются в map по их ID.
type MemStorage struct {
	metrics map[string]models.Metrics
}

// NewMemStorage создает новое хранилище метрик в памяти.
func NewMemStorage() handler.Storager {
	return &MemStorage{metrics: make(map[string]models.Metrics)}
}

// Save сохраняет метрику в хранилище.
// Для gauge просто перезаписывает значение.
// Для counter увеличивает значение счетчика, если метрика уже существует.
func (r *MemStorage) Save(metric models.Metrics) error {
	mType, ID, value, delta := metric.MType, metric.ID, metric.Value, metric.Delta

	var defaultDelta int64 = 0
	if delta == nil {
		delta = &defaultDelta
	}
	existMetric, ok := r.Get(mType, ID)

	if mType == models.Gauge {
		r.metrics[ID] = models.Metrics{ID: ID, MType: mType, Value: value}
		return nil
	}
	if mType == models.Counter {
		if !ok {
			r.metrics[ID] = models.Metrics{ID: ID, MType: mType, Delta: delta}
			return nil
		} else {
			delta := *existMetric.Delta + *delta
			existMetric.Delta = &delta
			r.metrics[ID] = existMetric
		}
	}
	return nil
}

// Get возвращает метрику по типу и ID.
// Если метрика не найдена или тип не совпадает — возвращает false.
func (r *MemStorage) Get(mType string, ID string) (models.Metrics, bool) {
	m, ok := r.metrics[ID]
	if !ok {
		return models.Metrics{}, false
	}
	if m.MType != mType {
		return models.Metrics{}, false
	}

	return m, true
}

// GetAll возвращает срез всех метрик, хранящихся в памяти.
func (r *MemStorage) GetAll() []models.Metrics {
	all := make([]models.Metrics, 0, len(r.metrics))
	for _, m := range r.metrics {
		fmt.Println(m)
		all = append(all, m)
	}
	return all
}
