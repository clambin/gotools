// Package metrics facilitates writing unit tests by allowing
// the test to read back the value set in a Prometheus metrics
package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
)

var (
	metrics = make(map[string]interface{})
)

// NewGauge returns a new Prometheus Gauge, created through promauto
func NewGauge(opts prometheus.GaugeOpts) prometheus.Gauge {
	metric := promauto.NewGauge(opts)
	metrics[opts.Name] = metric

	return metric
}

// NewGaugeVec returns a new Prometheus GaugeVec, created through promauto
func NewGaugeVec(opts prometheus.GaugeOpts, labels []string) *prometheus.GaugeVec {
	metric := promauto.NewGaugeVec(opts, labels)
	metrics[opts.Name] = metric

	return metric
}

// NewCounter returns a new Prometheus Counter, created through promauto
func NewCounter(opts prometheus.CounterOpts) prometheus.Counter {
	metric := promauto.NewCounter(opts)
	metrics[opts.Name] = metric

	return metric
}

// NewCounterVec returns a new Prometheus CounterVec, created through promauto
func NewCounterVec(opts prometheus.CounterOpts, labels []string) *prometheus.CounterVec {
	metric := promauto.NewCounterVec(opts, labels)
	metrics[opts.Name] = metric

	return metric
}

// LoadValue gets the last value reported so unit tests can verify the correct value was reported
func LoadValue(metricName string, labels ...string) (float64, error) {
	log.Debugf("%s(%s)", metricName, labels)
	if metric, ok := metrics[metricName]; ok {
		var m = dto.Metric{}
		switch metric.(type) {
		case prometheus.Gauge:
			gauge := metric.(prometheus.Gauge)
			_ = gauge.Write(&m)
			return m.Gauge.GetValue(), nil
		case *prometheus.GaugeVec:
			gaugevec := metric.(*prometheus.GaugeVec)
			log.Debug(gaugevec)
			_ = gaugevec.WithLabelValues(labels...).Write(&m)
			return m.Gauge.GetValue(), nil
		case prometheus.Counter:
			gauge := metric.(prometheus.Counter)
			_ = gauge.Write(&m)
			return m.Counter.GetValue(), nil
		case *prometheus.CounterVec:
			countervec := metric.(*prometheus.CounterVec)
			log.Debug(countervec)
			_ = countervec.WithLabelValues(labels...).Write(&m)
			return m.Counter.GetValue(), nil
		}
	}
	return 0, fmt.Errorf("could not find %s", metricName)
}
