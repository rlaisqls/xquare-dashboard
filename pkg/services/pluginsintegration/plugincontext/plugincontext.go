package plugincontext

import (
	"context"
	"encoding/json"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/xquare-dashboard/pkg/plugins"
	"github.com/xquare-dashboard/pkg/plugins/manager/store"
	"github.com/xquare-dashboard/pkg/services/datasources"
)

func ProvideService(pluginStore store.Service) *Provider {
	return &Provider{
		pluginStore: pluginStore,
	}
}

type Provider struct {
	pluginStore store.Service
}

// Get allows getting plugin context by its ID. If datasourceUID is not empty string
// then PluginContext.DataSourceInstanceSettings will be resolved and appended to
// returned context.
// Note: identity.Requester can be nil.
func (p *Provider) Get(ctx context.Context, pluginType datasources.DataSourceType, orgID int64) (backend.PluginContext, error) {
	plugin, exists := p.pluginStore.Plugin(ctx, string(pluginType))
	if !exists {
		return backend.PluginContext{}, plugins.ErrPluginNotRegistered
	}
	pCtx := backend.PluginContext{
		PluginID: plugin.ID,
	}
	return pCtx, nil
}

// GetWithDataSource allows getting plugin context by its ID and PluginContext.DataSourceInstanceSettings will be
// resolved and appended to the returned context.
// Note: *user.SignedInUser can be nil.
func (p *Provider) GetWithDataSource(ctx context.Context, pluginType datasources.DataSourceType, ds *datasources.DataSource) (backend.PluginContext, error) {
	plugin, exists := p.pluginStore.Plugin(ctx, string(pluginType))
	if !exists {
		return backend.PluginContext{}, plugins.ErrPluginNotRegistered
	}
	pCtx := backend.PluginContext{
		PluginID: plugin.ID,
	}
	name := string(ds.Type)
	pCtx.DataSourceInstanceSettings = &backend.DataSourceInstanceSettings{
		Type:     name,
		Name:     name,
		URL:      ds.URL,
		UID:      name,
		JSONData: json.RawMessage("{}"),
	}
	return pCtx, nil
}
