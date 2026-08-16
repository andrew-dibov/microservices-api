package registries

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewPromRegistry() *PromRegistry {
	p := prometheus.NewRegistry()

	p.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &PromRegistry{p: p}
}

func (r *PromRegistry) Metrics() http.Handler {
	return promhttp.HandlerFor(r.p, promhttp.HandlerOpts{
		Registry: r.p,
	})
}

func (r *PromRegistry) Add(collectors ...prometheus.Collector) {
	r.p.MustRegister(collectors...)
}

func (r *PromRegistry) Reg() *prometheus.Registry {
	return r.p
}
