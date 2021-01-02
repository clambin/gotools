package metrics_test

import (
	"github.com/clambin/gotools/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGauge(t *testing.T) {
	gauge := metrics.NewGauge(prometheus.GaugeOpts{
		Name: "test_gauge",
		Help: "Gauge test",
	})

	assert.NotNil(t, gauge)
	gauge.Set(1)

	value, err := metrics.LoadValue("test_gauge")
	assert.Nil(t, err)
	assert.Equal(t, float64(1), value)
}

func TestGaugeVec(t *testing.T) {
	gauge := metrics.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test_gaugevec",
		Help: "Gauge test",
	}, []string{"host"})

	assert.NotNil(t, gauge)
	gauge.WithLabelValues("host1").Set(2)

	value, err := metrics.LoadValue("test_gaugevec", "host1")
	assert.Nil(t, err)
	assert.Equal(t, float64(2), value)
}

func TestCounter(t *testing.T) {
	counter := metrics.NewCounter(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "Gauge test",
	})

	assert.NotNil(t, counter)
	counter.Add(10)

	value, err := metrics.LoadValue("test_counter")
	assert.Nil(t, err)
	assert.Equal(t, float64(10), value)
}

func TestCounterVec(t *testing.T) {
	counter := metrics.NewCounterVec(prometheus.CounterOpts{
		Name: "test_countervec",
		Help: "Gauge test",
	}, []string{"host"})

	assert.NotNil(t, counter)
	counter.WithLabelValues("host1").Add(20)

	value, err := metrics.LoadValue("test_countervec", "host1")
	assert.Nil(t, err)
	assert.Equal(t, float64(20), value)
}

func TestInvalid(t *testing.T) {
	_, err := metrics.LoadValue("invalid")
	assert.NotNil(t, err)
}
