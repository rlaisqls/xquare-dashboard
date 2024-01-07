package instrumentation

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	pluginParsingResponseDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "grafana",
		Name:      "loki_plugin_parse_response_duration_seconds",
		Help:      "Duration of Loki parsing the response in seconds",
		Buckets:   []float64{.001, 0.0025, .005, .0075, .01, .02, .03, .04, .05, .075, .1, .25, .5, 1, 5, 10, 25},
	}, []string{"status", "endpoint"})
)

const (
	EndpointQueryData = "queryData"
)

func UpdatePluginParsingResponseDurationSeconds(ctx context.Context, duration time.Duration, status string) {
	histogram := pluginParsingResponseDurationSeconds.WithLabelValues(status, EndpointQueryData)
	histogram.Observe(duration.Seconds())
}
