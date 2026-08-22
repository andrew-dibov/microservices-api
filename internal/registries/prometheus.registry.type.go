package registries

import "github.com/prometheus/client_golang/prometheus"

type PrometheusRegistry struct {
	prometheus *prometheus.Registry
}
