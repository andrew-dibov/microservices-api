package registries

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewPromRegistry() *PromRegistry {
	preg := prometheus.NewRegistry()

	preg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &PromRegistry{preg: preg}
}

func (r *PromRegistry) Metrics() http.Handler {
	return promhttp.HandlerFor(r.preg, promhttp.HandlerOpts{
		Registry: r.preg,
	})
}

func (r *PromRegistry) Add(collectors ...prometheus.Collector) {
	r.preg.MustRegister(collectors...)
}

func (r *PromRegistry) Reg() *prometheus.Registry {
	return r.preg
}
