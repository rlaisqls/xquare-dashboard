// Package api Grafana HTTP API.
//
// The Grafana backend exposes an HTTP API, the same API is used by the frontend to do
// everything from saving dashboards, creating users and updating data sources.
//
//	Schemes: http, https
//	BasePath: /api
//	Version: 0.0.1
//	Contact: Grafana Labs<hello@grafana.com> https://grafana.com
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- basic:
//	- api_key:
//
//	SecurityDefinitions:
//	basic:
//	 type: basic
//	api_key:
//	 type: apiKey
//	 name: Authorization
//	 in: header
//
// swagger:meta
package api

import (
	"github.com/xquare-dashboard/pkg/api/routing"
	"github.com/xquare-dashboard/pkg/infra/log"
	"github.com/xquare-dashboard/pkg/middleware/requestmeta"
)

var plog = log.New("api")

// registerRoutes registers all API HTTP routes.
func (hs *HTTPServer) registerRoutes() {
	r := hs.RouteRegister
	r.Group("/api", func(apiRoute routing.RouteRegister) {
		// metrics
		// DataSource w/ expressions
		apiRoute.Post("/ds/query", requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow), routing.Wrap(hs.QueryMetrics))
	})
}
