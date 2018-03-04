package route

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	domainCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "route",
			Subsystem: "digger",
			Name:      "domain",
			Help:      "Domain digged",
		},
		[]string{
			"domain",
			"type",
		},
	)
	clientCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "route",
			Subsystem: "digger",
			Name:      "client",
			Help:      "client queryed",
		},
		[]string{
			"ip",
		},
	)
)

func registerPrometheus() {
	prometheus.MustRegister(domainCounter)
	prometheus.MustRegister(clientCounter)
}
