package store

import (
	"context"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/xquare-dashboard/pkg/plugins/backendplugin"
	"github.com/xquare-dashboard/pkg/plugins/backendplugin/coreplugin"
	"github.com/xquare-dashboard/pkg/tsdb/loki"
	"github.com/xquare-dashboard/pkg/tsdb/prometheus"
	"sync"

	"github.com/xquare-dashboard/pkg/plugins"
)

type InMemory struct {
	store map[string]*plugins.Plugin
	mu    sync.RWMutex
}

func ProvideService(lk *loki.Service, pr *prometheus.Service) *InMemory {
	i := &InMemory{
		store: make(map[string]*plugins.Plugin),
	}

	lokiPlugin, _ := asBackendPlugin(lk)("loki")
	i.store["loki"] = &plugins.Plugin{
		ID:            "loki",
		Signature:     "loki",
		BackendClient: lokiPlugin,
	}

	prometheusPlugin, _ := asBackendPlugin(pr)("prometheus")
	i.store["prometheus"] = &plugins.Plugin{
		ID:            "prometheus",
		Signature:     "prometheus",
		BackendClient: prometheusPlugin,
	}
	return i
}

func asBackendPlugin(svc any) backendplugin.PluginFactoryFunc {

	opts := backend.ServeOpts{}
	if queryHandler, ok := svc.(backend.QueryDataHandler); ok {
		opts.QueryDataHandler = queryHandler
	}
	if resourceHandler, ok := svc.(backend.CallResourceHandler); ok {
		opts.CallResourceHandler = resourceHandler
	}
	if streamHandler, ok := svc.(backend.StreamHandler); ok {
		opts.StreamHandler = streamHandler
	}
	if healthHandler, ok := svc.(backend.CheckHealthHandler); ok {
		opts.CheckHealthHandler = healthHandler
	}

	if opts.QueryDataHandler != nil || opts.CallResourceHandler != nil ||
		opts.CheckHealthHandler != nil || opts.StreamHandler != nil {
		return coreplugin.New(opts)
	}

	return nil
}

func (i *InMemory) Plugin(_ context.Context, pluginID string) (*plugins.Plugin, bool) {
	return i.plugin(pluginID)
}

func (i *InMemory) Plugins(_ context.Context) []*plugins.Plugin {
	i.mu.RLock()
	defer i.mu.RUnlock()

	res := make([]*plugins.Plugin, 0, len(i.store))
	for _, p := range i.store {
		res = append(res, p)
	}

	return res
}

func (i *InMemory) plugin(pluginID string) (*plugins.Plugin, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	p, exists := i.store[pluginID]

	if !exists {
		return nil, false
	}
	return p, true
}

func (i *InMemory) isRegistered(pluginID string) bool {
	p, exists := i.plugin(pluginID)

	// This may have matched based on an alias
	return exists && p.ID == pluginID
}
