package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Metrics struct {
	activeUsersGauge prometheus.Gauge
	Registry         *prometheus.Registry
}

func Register() *Metrics {
	activeUsersGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users_total",
			Help: "Number of current active users",
		},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		activeUsersGauge,
	)

	return &Metrics{
		activeUsersGauge: activeUsersGauge,
		Registry:         registry,
	}
}

func (m *Metrics) Increase() {
	m.activeUsersGauge.Inc()
}

func (m *Metrics) Decrease() {
	m.activeUsersGauge.Dec()
}
