package registries

import "github.com/prometheus/client_golang/prometheus"

type PromRegistry struct {
	preg *prometheus.Registry
}
