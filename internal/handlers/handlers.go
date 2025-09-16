// Package handler реализует HTTP-обработчики для работы с метриками.
// Обработчики принимают, сохраняют и возвращают метрики через различные HTTP-эндпоинты.

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexkozopolianski/go-metrics-tpl/internal/models"
	"github.com/go-chi/chi/v5"
)

// Handler — структура, инкапсулирующая хранилище метрик.
type Handler struct {
	storage Storager // Интерфейс хранилища метрик
}

// Storager — интерфейс для абстракции хранилища метрик.
type Storager interface {
	Save(metric models.Metrics) error
	Get(mType, id string) (models.Metrics, bool)
	GetAll() []models.Metrics
}

// NewHandler создает новый экземпляр Handler с заданным хранилищем.
func NewHandler(storage Storager) *Handler {
	return &Handler{storage: storage}
}

// Update — HTTP-обработчик для обновления метрики через параметры URL.
// Поддерживает типы gauge и counter. Валидирует входные данные.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "type")
	id := chi.URLParam(r, "id")
	value := chi.URLParam(r, "value")

	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if mType == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if mType != models.Counter && mType != models.Gauge {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	metric := models.Metrics{
		ID:    id,
		MType: mType,
	}

	if mType == models.Counter {
		parseInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric.Delta = &parseInt
	} else {
		parseFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metric.Value = &parseFloat
	}

	err := h.storage.Save(metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateJSON — HTTP-обработчик для обновления метрики через JSON в теле запроса.
// Возвращает обновлённую метрику в ответе.
func (h *Handler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	requestMetric := models.Metrics{}
	err := json.NewDecoder(r.Body).Decode(&requestMetric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var responseMetric models.Metrics

	err = h.storage.Save(requestMetric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metric, ok := h.storage.Get(requestMetric.MType, requestMetric.ID)
	if !ok {
		responseMetric = requestMetric
	} else {
		responseMetric = metric
	}

	res, err := json.Marshal(responseMetric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Value — HTTP-обработчик для получения значения метрики по типу и id через URL.
// Возвращает значение метрики в формате JSON.
func (h *Handler) Value(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "type")
	id := chi.URLParam(r, "id")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	m, ok := h.storage.Get(mType, id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var bytes []byte
	var err error
	if m.MType == models.Counter {
		bytes, err = json.Marshal(m.Delta)
	} else {
		bytes, err = json.Marshal(m.Value)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ValueJSON — HTTP-обработчик для получения метрики через JSON-запрос.
// Возвращает всю структуру метрики в формате JSON.
func (h *Handler) ValueJSON(w http.ResponseWriter, r *http.Request) {
	metric := models.Metrics{}

	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")

	m, ok := h.storage.Get(metric.MType, metric.ID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// All — HTTP-обработчик для получения всех метрик.
// Возвращает список всех метрик в формате JSON.
func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	m := h.storage.GetAll()

	bytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
