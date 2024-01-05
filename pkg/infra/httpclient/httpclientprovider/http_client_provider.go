package httpclientprovider

import (
	"net/http"
	"time"

	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/mwitkow/go-conntrack"

	"github.com/xquare-dashboard/pkg/infra/log"
	"github.com/xquare-dashboard/pkg/infra/metrics/metricutil"
	"github.com/xquare-dashboard/pkg/infra/tracing"
	"github.com/xquare-dashboard/pkg/services/validations"
)

var newProviderFunc = sdkhttpclient.NewProvider

// New creates a new HTTP client provider with pre-configured middlewares.
func New(validator validations.PluginRequestValidator, tracer tracing.Tracer) *sdkhttpclient.Provider {
	logger := log.New("httpclient")

	middlewares := []sdkhttpclient.Middleware{
		TracingMiddleware(logger, tracer),
		DataSourceMetricsMiddleware(),
		sdkhttpclient.ContextualMiddleware(),
		sdkhttpclient.BasicAuthenticationMiddleware(),
		sdkhttpclient.CustomHeadersMiddleware(),
		ResponseLimitMiddleware(100),
		RedirectLimitMiddleware(validator),
	}

	setDefaultTimeoutOptions()

	return newProviderFunc(sdkhttpclient.ProviderOptions{
		Middlewares: middlewares,
		ConfigureTransport: func(opts sdkhttpclient.Options, transport *http.Transport) {
			datasourceName, exists := opts.Labels["datasource_name"]
			if !exists {
				return
			}
			datasourceLabelName, err := metricutil.SanitizeLabelName(datasourceName)
			if err != nil {
				return
			}

			newConntrackRoundTripper(datasourceLabelName, transport)
		},
	})
}

// newConntrackRoundTripper takes a http.DefaultTransport and adds the Conntrack Dialer
// so we can instrument outbound connections
func newConntrackRoundTripper(name string, transport *http.Transport) *http.Transport {
	transport.DialContext = conntrack.NewDialContextFunc(
		conntrack.DialWithName(name),
		conntrack.DialWithDialContextFunc(transport.DialContext),
	)
	return transport
}

// setDefaultTimeoutOptions overrides the default timeout options for the SDK.
//
// Note: Not optimal changing global state, but hard to not do in this case.
func setDefaultTimeoutOptions() {
	sdkhttpclient.DefaultTimeoutOptions = sdkhttpclient.TimeoutOptions{
		Timeout:               time.Duration(10) * time.Second,
		DialTimeout:           time.Duration(30) * time.Second,
		KeepAlive:             time.Duration(30) * time.Second,
		TLSHandshakeTimeout:   time.Duration(10) * time.Second,
		ExpectContinueTimeout: time.Duration(1) * time.Second,
		MaxConnsPerHost:       10,
		MaxIdleConns:          10,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       time.Duration(90) * time.Second,
	}
}
