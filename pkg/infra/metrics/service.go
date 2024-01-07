package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func ProvideRegisterer() prometheus.Registerer {
	return prometheus.DefaultRegisterer
}

func ProvideGatherer() prometheus.Gatherer {
	return prometheus.DefaultGatherer
}
