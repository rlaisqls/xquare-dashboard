package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/xquare-dashboard/pkg/api/routing"
	"github.com/xquare-dashboard/pkg/infra/log"
	"github.com/xquare-dashboard/pkg/infra/tracing"
	"github.com/xquare-dashboard/pkg/web"
)

type HTTPServer struct {
	log              log.Logger
	web              *web.Mux
	context          context.Context
	httpSrv          *http.Server
	middlewares      []web.Handler
	namedMiddlewares []routing.RegisterNamedMiddleware

	RouteRegister routing.RouteRegister
	Listener      net.Listener
	tracer        tracing.Tracer

	promRegister prometheus.Registerer
	promGatherer prometheus.Gatherer
}

type ServerOptions struct {
	Listener net.Listener
}

func ProvideHTTPServer(
	opts ServerOptions,
	tracer tracing.Tracer,
	promGatherer prometheus.Gatherer,
	promRegister prometheus.Registerer,
) (*HTTPServer, error) {
	m := web.New()

	hs := &HTTPServer{
		tracer:       tracer,
		log:          log.New("http.server"),
		web:          m,
		Listener:     opts.Listener,
		promRegister: promRegister,
		promGatherer: promGatherer,
	}
	if hs.Listener != nil {
		hs.log.Debug("Using provided listener")
	}

	return hs, nil
}

func (hs *HTTPServer) AddMiddleware(middleware web.Handler) {
	hs.middlewares = append(hs.middlewares, middleware)
}

func (hs *HTTPServer) AddNamedMiddleware(middleware routing.RegisterNamedMiddleware) {
	hs.namedMiddlewares = append(hs.namedMiddlewares, middleware)
}

func (hs *HTTPServer) applyRoutes() {
	// start with middlewares & static routes
	hs.addMiddlewaresAndStaticRoutes()
	// then add view routes & api routes
	hs.RouteRegister.Register(hs.web, hs.namedMiddlewares...)
}

func (hs *HTTPServer) addMiddlewaresAndStaticRoutes() {
	m := hs.web

	// These endpoints are used for monitoring the Grafana instance
	// and should not be redirected or rejected.
	m.Use(hs.healthzHandler)
	m.Use(hs.metricsEndpoint)

	for _, mw := range hs.middlewares {
		m.Use(mw)
	}
}

func (hs *HTTPServer) Run(ctx context.Context) error {
	hs.context = ctx

	hs.applyRoutes()

	// Remove any square brackets enclosing IPv6 addresses, a format we support for backwards compatibility
	host := strings.TrimSuffix(strings.TrimPrefix("dashboard.xquare.app", "["), "]")
	hs.httpSrv = &http.Server{
		Addr:        net.JoinHostPort(host, "9090"),
		Handler:     hs.web,
		ReadTimeout: 30,
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

	if err := hs.httpSrv.ServeTLS(listener, "", ""); err != nil {
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
	if hs.Listener != nil {
		return hs.Listener, nil
	}

	listener, err := net.Listen("tcp", hs.httpSrv.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to open listener on address %s: %w", hs.httpSrv.Addr, err)
	}
	return listener, nil
}

func (hs *HTTPServer) metricsEndpoint(ctx *web.Context) {

	if ctx.Req.Method != http.MethodGet || ctx.Req.URL.Path != "/metrics" {
		return
	}

	promhttp.
		HandlerFor(hs.promGatherer, promhttp.HandlerOpts{EnableOpenMetrics: true}).
		ServeHTTP(ctx.Resp, ctx.Req)
}

// healthzHandler always return 200 - Ok if Grafana's web server is running
func (hs *HTTPServer) healthzHandler(ctx *web.Context) {
	notHeadOrGet := ctx.Req.Method != http.MethodGet && ctx.Req.Method != http.MethodHead
	if notHeadOrGet || ctx.Req.URL.Path != "/healthz" {
		return
	}

	ctx.Resp.WriteHeader(http.StatusOK)
	if _, err := ctx.Resp.Write([]byte("Ok")); err != nil {
		hs.log.Error("could not write to response", "err", err)
	}
}
