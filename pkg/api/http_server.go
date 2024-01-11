package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xquare-dashboard/pkg/api/routing"
	"github.com/xquare-dashboard/pkg/middleware"
	"github.com/xquare-dashboard/pkg/plugins"
	"github.com/xquare-dashboard/pkg/plugins/manager/store"
	"github.com/xquare-dashboard/pkg/services/pluginsintegration/plugincontext"
	"net"
	"net/http"
	"path"
	"sync"

	httpstatic "github.com/xquare-dashboard/pkg/api/static"
	"github.com/xquare-dashboard/pkg/components/simplejson"
	"github.com/xquare-dashboard/pkg/infra/log"
	"github.com/xquare-dashboard/pkg/services/contexthandler"
	"github.com/xquare-dashboard/pkg/services/query"
	"github.com/xquare-dashboard/pkg/web"
)

type HTTPServer struct {
	log              log.Logger
	web              *web.Mux
	context          context.Context
	httpSrv          *http.Server
	middlewares      []web.Handler
	RouteRegister    routing.RouteRegister
	namedMiddlewares []routing.RegisterNamedMiddleware
	pluginStore      store.Service
	pluginClient     plugins.Client

	ContextHandler   *contexthandler.ContextHandler
	pCtxProvider     *plugincontext.Provider
	queryDataService query.Service
	promRegister     prometheus.Registerer
	promGatherer     prometheus.Gatherer
}

func ProvideHTTPServer(
	contextHandler *contexthandler.ContextHandler, queryDataService query.Service,
	promGatherer prometheus.Gatherer, promRegister prometheus.Registerer, pluginClient plugins.Client,
	routeRegister routing.RouteRegister, pluginStore store.Service, pCtxProvider *plugincontext.Provider,
) (*HTTPServer, error) {
	m := web.New()
	hs := &HTTPServer{
		ContextHandler:   contextHandler,
		log:              log.New("http.server"),
		web:              m,
		queryDataService: queryDataService,
		pluginClient:     pluginClient,
		promRegister:     promRegister,
		promGatherer:     promGatherer,
		RouteRegister:    routeRegister,
		pluginStore:      pluginStore,
		pCtxProvider:     pCtxProvider,
	}
	hs.registerRoutes()
	return hs, nil
}

func (hs *HTTPServer) AddMiddleware(middleware web.Handler) {
	hs.middlewares = append(hs.middlewares, middleware)
}

func (hs *HTTPServer) Run(ctx context.Context) error {

	hs.context = ctx
	hs.applyRoutes()

	// Remove any square brackets enclosing IPv6 addresses, a format we support for backwards compatibility
	hs.httpSrv = &http.Server{
		Addr:        ":9090",
		Handler:     hs.web,
		ReadTimeout: 10000,
	}
	if err := hs.configureHttp2(); err != nil {
		return err
	}

	listener, err := hs.getListener()
	if err != nil {
		return err
	}

	hs.log.Info("HTTP Server Listen", "address", listener.Addr().String())

	var wg sync.WaitGroup
	wg.Add(1)

	// handle http shutdown on server context done
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := hs.httpSrv.Shutdown(context.Background()); err != nil {
			hs.log.Error("Failed to shutdown server", "error", err)
		}
	}()

	if err := hs.httpSrv.Serve(listener); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			hs.log.Debug("server was shutdown gracefully")
			return nil
		}
		return err
	}
	wg.Wait()

	return nil
}

func (hs *HTTPServer) getListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", hs.httpSrv.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to open listener on address %s: %w", hs.httpSrv.Addr, err)
	}
	return listener, nil
}

func (hs *HTTPServer) configureHttp2() error {
	return nil
}

func (hs *HTTPServer) applyRoutes() {
	// start with middlewares & static routes
	hs.addMiddlewaresAndStaticRoutes()
	// then add view routes & api routes
	hs.RouteRegister.Register(hs.web, hs.namedMiddlewares...)
	// lastly not found route
	hs.web.NotFound(middleware.ProvideRouteOperationName("notfound"), hs.NotFoundHandler)
}

func (hs *HTTPServer) addMiddlewaresAndStaticRoutes() {
	m := hs.web

	// These endpoints are used for monitoring the Grafana instance
	// and should not be redirected or rejected.
	m.Use(hs.apiHealthHandler)
	m.Use(hs.metricsEndpoint)

	m.UseMiddleware(hs.ContextHandler.Middleware)

	for _, mw := range hs.middlewares {
		m.Use(mw)
	}
}

func (hs *HTTPServer) metricsEndpoint(ctx *web.Context) {

	if ctx.Req.Method != http.MethodGet || ctx.Req.URL.Path != "/metrics" {
		return
	}

	promhttp.
		HandlerFor(hs.promGatherer, promhttp.HandlerOpts{EnableOpenMetrics: true}).
		ServeHTTP(ctx.Resp, ctx.Req)
}

// apiHealthHandler will return ok if web server is running.
func (hs *HTTPServer) apiHealthHandler(ctx *web.Context) {
	notHeadOrGet := ctx.Req.Method != http.MethodGet && ctx.Req.Method != http.MethodHead
	if notHeadOrGet || ctx.Req.URL.Path != "/api/health" {
		return
	}

	data := simplejson.New()
	data.Set("database", "ok")
	dataBytes, err := data.EncodePretty()
	if err != nil {
		hs.log.Error("Failed to encode data", "err", err)
		return
	}

	if _, err := ctx.Resp.Write(dataBytes); err != nil {
		hs.log.Error("Failed to write to response", "err", err)
	}
}

func (hs *HTTPServer) mapStatic(m *web.Mux, rootDir string, dir string, prefix string, exclude ...string) {
	headers := func(c *web.Context) {
		c.Resp.Header().Set("Cache-Control", "public, max-age=3600")
	}

	m.Use(httpstatic.Static(
		path.Join(rootDir, dir),
		httpstatic.StaticOptions{
			SkipLogging: true,
			Prefix:      prefix,
			AddHeaders:  headers,
			Exclude:     exclude,
		},
	))
}
