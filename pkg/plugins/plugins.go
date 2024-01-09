package plugins

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/xquare-dashboard/pkg/plugins/backendplugin"
)

var (
	ErrFileNotExist              = errors.New("file does not exist")
	ErrPluginFileRead            = errors.New("file could not be read")
	ErrUninstallInvalidPluginDir = errors.New("cannot recognize as plugin folder")
	ErrInvalidPluginJSON         = errors.New("did not find valid type or id properties in plugin.json")
	ErrUnsupportedAlias          = errors.New("can not set alias in plugin.json")
)

type Plugin struct {
	ID            string
	Signature     string
	BackendClient backendplugin.BackendPlugin
	mu            sync.Mutex
}

func (p *Plugin) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.BackendClient == nil {
		return fmt.Errorf("could not start plugin %s as no plugin BackendClient exists", p.ID)
	}

	return p.BackendClient.Start(ctx)
}

func (p *Plugin) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.BackendClient == nil {
		return nil
	}

	return p.BackendClient.Stop(ctx)
}

func (p *Plugin) IsManaged() bool {
	if p.BackendClient != nil {
		return p.BackendClient.IsManaged()
	}
	return false
}

func (p *Plugin) Decommission() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.BackendClient != nil {
		return p.BackendClient.Decommission()
	}
	return nil
}

func (p *Plugin) IsDecommissioned() bool {
	if p.BackendClient != nil {
		return p.BackendClient.IsDecommissioned()
	}
	return false
}

func (p *Plugin) Exited() bool {
	if p.BackendClient != nil {
		return p.BackendClient.Exited()
	}
	return false
}

func (p *Plugin) Target() backendplugin.Target {
	if p.BackendClient == nil {
		return backendplugin.TargetUnknown
	}
	return p.BackendClient.Target()
}

func (p *Plugin) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	pluginClient, ok := p.Client()
	if !ok {
		return nil, ErrPluginUnavailable
	}
	return pluginClient.QueryData(ctx, req)
}

func (p *Plugin) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	pluginClient, ok := p.Client()
	if !ok {
		return ErrPluginUnavailable
	}
	return pluginClient.CallResource(ctx, req, sender)
}

func (p *Plugin) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	pluginClient, ok := p.Client()
	if !ok {
		return nil, ErrPluginUnavailable
	}
	return pluginClient.CheckHealth(ctx, req)
}

func (p *Plugin) CollectMetrics(ctx context.Context, req *backend.CollectMetricsRequest) (*backend.CollectMetricsResult, error) {
	pluginClient, ok := p.Client()
	if !ok {
		return nil, ErrPluginUnavailable
	}
	return pluginClient.CollectMetrics(ctx, req)
}

func (p *Plugin) SubscribeStream(ctx context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	pluginClient, ok := p.Client()
	if !ok {
		return nil, ErrPluginUnavailable
	}
	return pluginClient.SubscribeStream(ctx, req)
}

func (p *Plugin) PublishStream(ctx context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	pluginClient, ok := p.Client()
	if !ok {
		return nil, ErrPluginUnavailable
	}
	return pluginClient.PublishStream(ctx, req)
}

func (p *Plugin) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	pluginClient, ok := p.Client()
	if !ok {
		return ErrPluginUnavailable
	}
	return pluginClient.RunStream(ctx, req, sender)
}

func (p *Plugin) RegisterClient(c backendplugin.BackendPlugin) {
	p.BackendClient = c
}

func (p *Plugin) Client() (PluginClient, bool) {
	if p.BackendClient != nil {
		return p.BackendClient, true
	}
	return nil, false
}

type PluginClient interface {
	backend.QueryDataHandler
	backend.CollectMetricsHandler
	backend.CheckHealthHandler
	backend.CallResourceHandler
	backend.StreamHandler
}

var PluginTypes = []Type{
	TypeDataSource,
}

type Type string

const (
	TypeDataSource Type = "datasource"
)

func (pt Type) IsValid() bool {
	switch pt {
	case TypeDataSource:
		return true
	}
	return false
}
