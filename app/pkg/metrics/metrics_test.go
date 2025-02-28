package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegisterAndGetHistogram(t *testing.T) {
	m := Init()
	key := "test_histogram"
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: key,
			Help: "Test histogram",
		}, []string{"label"},
	)

	err := m.RegisterHistogram(key, histogram)
	require.NoError(t, err)

	retrieved, err := m.GetHistogram(key)
	require.NoError(t, err)
	assert.Equal(t, histogram, retrieved)
}

func TestRegisterAndGetCounter(t *testing.T) {
	m := Init()
	key := "test_counter"

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: key,
			Help: "Test counter",
		}, []string{"label"},
	)

	err := m.RegisterCounter(key, counter)
	require.NoError(t, err)

	retrieved, err := m.GetCounter(key)
	require.NoError(t, err)
	assert.Equal(t, counter, retrieved)
}

func TestRegisterHistogramDuplicate(t *testing.T) {
	m := Init()
	key := "test_histogram_dup"

	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: key,
			Help: "Test histogram",
		}, []string{"label"},
	)

	err := m.RegisterHistogram(key, histogram)
	require.NoError(t, err)

	// Попытка повторной регистрации
	err = m.RegisterHistogram(key, histogram)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestRegisterCounterDuplicate(t *testing.T) {
	m := Init()
	key := "test_counter_dup"

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: key,
			Help: "Test counter",
		}, []string{"label"},
	)

	err := m.RegisterCounter(key, counter)
	require.NoError(t, err)

	// Попытка повторной регистрации
	err = m.RegisterCounter(key, counter)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestGetNonExistentMetrics(t *testing.T) {
	m := Init()

	_, err := m.GetHistogram("non_existent_histogram")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	_, err = m.GetCounter("non_existent_counter")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
