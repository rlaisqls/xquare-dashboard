package backendplugin

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// BackendPlugin is the backend plugin interface.
type BackendPlugin interface {
	PluginID() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	IsManaged() bool
	Exited() bool
	Decommission() error
	IsDecommissioned() bool
	Target() Target
	backend.CollectMetricsHandler
	backend.CheckHealthHandler
	backend.QueryDataHandler
	backend.CallResourceHandler
	backend.StreamHandler
}

type Target string

const (
	TargetNone     Target = "none"
	TargetUnknown  Target = "unknown"
	TargetInMemory Target = "in_memory"
	TargetLocal    Target = "local"
)
