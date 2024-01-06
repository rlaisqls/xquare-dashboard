//go:build wireinject
// +build wireinject

// This file should contain wire sets used by both OSS and Enterprise builds.
// Use wireext_oss.go and wireext_enterprise.go for sets that are specific to
// the respective builds.
package server

import (
	"github.com/google/wire"
	"github.com/xquare-dashboard/pkg/api"
	"github.com/xquare-dashboard/pkg/expr"

	"github.com/xquare-dashboard/pkg/services/contexthandler"
	"github.com/xquare-dashboard/pkg/services/query"
	"github.com/xquare-dashboard/pkg/tsdb/legacydata"
	legacydataservice "github.com/xquare-dashboard/pkg/tsdb/legacydata/service"
	"github.com/xquare-dashboard/pkg/tsdb/loki"
	"github.com/xquare-dashboard/pkg/tsdb/prometheus"
)

var wireSet = wire.NewSet(
	legacydataservice.ProvideService,
	wire.Bind(new(legacydata.RequestHandler), new(*legacydataservice.Service)),
	New,
	api.ProvideHTTPServer,
	query.ProvideService,
	wire.Bind(new(query.Service), new(*query.ServiceImpl)),
	contexthandler.ProvideService,
	loki.ProvideService,
	prometheus.ProvideService,
	expr.ProvideService,
)

func Initialize() (*Server, error) {
	wire.Build(wireSet)
	return &Server{}, nil
}
