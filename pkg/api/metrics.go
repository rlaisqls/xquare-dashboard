package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/xquare-dashboard/pkg/api/dtos"
	"github.com/xquare-dashboard/pkg/api/response"
	"github.com/xquare-dashboard/pkg/middleware/requestmeta"
	"github.com/xquare-dashboard/pkg/plugins"
	"github.com/xquare-dashboard/pkg/plugins/httpresponsesender"
	contextmodel "github.com/xquare-dashboard/pkg/services/contexthandler/model"
	"github.com/xquare-dashboard/pkg/services/datasources"
	"github.com/xquare-dashboard/pkg/web"
	"io"
	"net/http"
	"net/url"
)

// QueryMetrics returns query metrics.
// swagger:route POST /ds/query ds queryMetricsWithExpressions
//
// DataSource query metrics with expressions.
//
// If you are running Grafana Enterprise and have Fine-grained access control enabled
// you need to have a permission with action: `datasources:query`.
//
// Responses:
// 200: queryMetricsWithExpressionsRespons
// 207: queryMetricsWithExpressionsRespons
// 401: unauthorisedError
// 400: badRequestError
// 403: forbiddenError
// 500: internalServerError
func (hs *HTTPServer) QueryMetrics(c *contextmodel.ReqContext) response.Response {
	reqDTO := dtos.MetricRequest{}
	if err := web.Bind(c.Req, &reqDTO); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}

	resp, err := hs.queryDataService.QueryData(c.Req.Context(), reqDTO)
	if err != nil {
		return hs.handleQueryMetricsError(err)
	}
	return hs.toJsonStreamingResponse(c.Req.Context(), resp)
}

func (hs *HTTPServer) handleQueryMetricsError(err error) *response.NormalResponse {
	return response.ErrOrFallback(http.StatusInternalServerError, "Query data error", err)
}

func (hs *HTTPServer) toJsonStreamingResponse(ctx context.Context, qdr *backend.QueryDataResponse) response.Response {
	statusWhenError := http.StatusBadRequest

	statusCode := http.StatusOK
	for _, res := range qdr.Responses {
		if res.Error != nil {
			statusCode = statusWhenError
		}
	}

	if statusCode == statusWhenError {
		// an error in the response we treat as downstream.
		requestmeta.WithDownstreamStatusSource(ctx)
	}

	return response.JSONStreaming(statusCode, qdr)
}

// Fetch data source resources.
//
// Responses:
// 200: okResponse
// 400: badRequestError
// 401: unauthorisedError
// 403: forbiddenError
// 404: notFoundError
// 500: internalServerError
func (hs *HTTPServer) CallDatasourceResourceWithUID(c *contextmodel.ReqContext) {
	dsUID := web.Params(c.Req)[":uid"]
	ds, _ := datasources.GetDataSource(datasources.DataSourceType(dsUID))

	plugin, exists := hs.pluginStore.Plugin(c.Req.Context(), string(ds.Type))
	if !exists {
		c.JsonApiErr(http.StatusInternalServerError, "Unable to find datasource plugin", nil)
		return
	}

	hs.callPluginResourceWithDataSource(c, plugin.ID, ds)
}

func (hs *HTTPServer) callPluginResourceWithDataSource(c *contextmodel.ReqContext, pluginID string, ds *datasources.DataSource) {
	pCtx, err := hs.pCtxProvider.GetWithDataSource(c.Req.Context(), datasources.DataSourceType(pluginID), ds)
	if err != nil {
		if errors.Is(err, plugins.ErrPluginNotRegistered) {
			c.JsonApiErr(404, "Plugin not found", nil)
			return
		}
		c.JsonApiErr(500, "Failed to get plugin settings", err)
		return
	}

	req, err := hs.pluginResourceRequest(c)
	if err != nil {
		c.JsonApiErr(http.StatusBadRequest, "Failed for create plugin resource request", err)
		return
	}

	if err = hs.makePluginResourceRequest(c.Resp, req, pCtx); err != nil {
		handleCallResourceError(err, c)
		return
	}

	requestmeta.WithStatusSource(c.Req.Context(), c.Resp.Status())
}

func (hs *HTTPServer) makePluginResourceRequest(w http.ResponseWriter, req *http.Request, pCtx backend.PluginContext) error {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	crReq := &backend.CallResourceRequest{
		PluginContext: pCtx,
		Path:          req.URL.Path,
		Method:        req.Method,
		URL:           req.URL.String(),
		Headers:       req.Header,
		Body:          body,
	}

	httpSender := httpresponsesender.New(w)
	return hs.pluginClient.CallResource(req.Context(), crReq, httpSender)
}

func (hs *HTTPServer) pluginResourceRequest(c *contextmodel.ReqContext) (*http.Request, error) {
	clonedReq := c.Req.Clone(c.Req.Context())
	rawURL := web.Params(c.Req)["*"]
	if clonedReq.URL.RawQuery != "" {
		rawURL += "?" + clonedReq.URL.RawQuery
	}
	urlPath, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	clonedReq.URL = urlPath

	return clonedReq, nil
}

func handleCallResourceError(err error, reqCtx *contextmodel.ReqContext) {
	resp := response.ErrOrFallback(http.StatusInternalServerError, "Failed to call resource", err)
	resp.WriteTo(reqCtx)
}
