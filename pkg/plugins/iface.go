package plugins

import (
	"context"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type PluginSource interface {
	PluginURIs(ctx context.Context) []string
	DefaultSignature(ctx context.Context) (string, bool)
}

// Client is used to communicate with backend plugin implementations.
type Client interface {
	backend.QueryDataHandler
	backend.CheckHealthHandler
	backend.StreamHandler
	backend.CallResourceHandler
	backend.CollectMetricsHandler
}
