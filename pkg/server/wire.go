//go:build wireinject
// +build wireinject

// This file should contain wire sets used by both OSS and Enterprise builds.
// Use wireext_oss.go and wireext_enterprise.go for sets that are specific to
// the respective builds.
package server

import (
	"github.com/google/wire"
	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/xquare-dashboard/pkg/api"
	"github.com/xquare-dashboard/pkg/api/routing"
	"github.com/xquare-dashboard/pkg/infra/httpclient"
	"github.com/xquare-dashboard/pkg/infra/httpclient/httpclientprovider"
	"github.com/xquare-dashboard/pkg/infra/metrics"
	"github.com/xquare-dashboard/pkg/plugins"
	"github.com/xquare-dashboard/pkg/plugins/manager/client"
	"github.com/xquare-dashboard/pkg/plugins/manager/store"
	"github.com/xquare-dashboard/pkg/registry"
	"github.com/xquare-dashboard/pkg/registry/backgroundsvcs"
	"github.com/xquare-dashboard/pkg/services/pluginsintegration/plugincontext"

	"github.com/xquare-dashboard/pkg/services/contexthandler"
	"github.com/xquare-dashboard/pkg/services/query"
	"github.com/xquare-dashboard/pkg/tsdb/loki"
	"github.com/xquare-dashboard/pkg/tsdb/prometheus"
)

var wireSet = wire.NewSet(
	New,
	api.ProvideHTTPServer,
	query.ProvideService,
	wire.Bind(new(query.Service), new(*query.ServiceImpl)),
	routing.ProvideRegister,
	wire.Bind(new(routing.RouteRegister), new(*routing.RouteRegisterImpl)),
	httpclientprovider.New,
	wire.Bind(new(httpclient.Provider), new(*sdkhttpclient.Provider)),
	contexthandler.ProvideService,
	loki.ProvideService,
	prometheus.ProvideService,
	store.ProvideService,
	wire.Bind(new(store.Service), new(*store.InMemory)),
	backgroundsvcs.ProvideBackgroundServiceRegistry,
	wire.Bind(new(registry.BackgroundServiceRegistry), new(*backgroundsvcs.BackgroundServiceRegistry)),
	plugincontext.ProvideService,
	client.ProvideService,
	wire.Bind(new(plugins.Client), new(*client.Service)),
	metrics.ProvideRegisterer,
	metrics.ProvideGatherer,
)

func Initialize() (*Server, error) {
	wire.Build(wireSet)
	return &Server{}, nil
}
