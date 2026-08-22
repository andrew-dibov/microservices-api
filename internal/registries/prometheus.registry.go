package registries

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewPrometheusRegistry() *PrometheusRegistry {
	prometheus := prometheus.NewRegistry()

	prometheus.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &PrometheusRegistry{prometheus: prometheus}
}

func (registry *PrometheusRegistry) Metrics() http.Handler {
	return promhttp.HandlerFor(registry.prometheus, promhttp.HandlerOpts{Registry: registry.prometheus})
}
