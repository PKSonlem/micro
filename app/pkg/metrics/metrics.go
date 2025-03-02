package metrics

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetrics struct {
	mu         sync.RWMutex
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
}

func Init() *PrometheusMetrics {
	return &PrometheusMetrics{
		counters:   make(map[string]*prometheus.CounterVec),
		histograms: make(map[string]*prometheus.HistogramVec),
	}
}

func (m *PrometheusMetrics) RegisterHistogram(key string, metric *prometheus.HistogramVec) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.histograms[key]; exists {
		return fmt.Errorf("histogramVec already exists with key: %s", key)
	}

	m.histograms[key] = metric
	prometheus.MustRegister(metric)

	return nil
}

func (m *PrometheusMetrics) RegisterCounter(key string, metric *prometheus.CounterVec) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.counters[key]; exists {
		return fmt.Errorf("counterVec already exists with key: %s", key)
	}

	m.counters[key] = metric
	prometheus.MustRegister(metric)

	return nil
}

func (m *PrometheusMetrics) GetHistogram(key string) (*prometheus.HistogramVec, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	histogram, exists := m.histograms[key]
	if !exists {
		return nil, fmt.Errorf("histogramVec not found with key: %s", key)
	}

	return histogram, nil
}

func (m *PrometheusMetrics) GetCounter(key string) (*prometheus.CounterVec, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	counter, exists := m.counters[key]
	if !exists {
		return nil, fmt.Errorf("counterVec not found with key: %s", key)
	}

	return counter, nil
}
