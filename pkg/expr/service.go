package expr

import (
	"context"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xquare-dashboard/pkg/plugins"
	"github.com/xquare-dashboard/pkg/services/datasources"
	"github.com/xquare-dashboard/pkg/services/pluginsintegration/plugincontext"
	"time"
)

// Service is service representation for expression handling.
type Service struct {
	metrics      *metrics
	dataService  backend.QueryDataHandler
	pCtxProvider pluginContextProvider
}

type pluginContextProvider interface {
	Get(ctx context.Context, pluginType datasources.DataSourceType, orgID int64) (backend.PluginContext, error)
	GetWithDataSource(ctx context.Context, pluginType datasources.DataSourceType, ds *datasources.DataSource) (backend.PluginContext, error)
}

func ProvideService(registerer prometheus.Registerer, pluginClient plugins.Client, pCtxProvider *plugincontext.Provider) *Service {
	return &Service{
		dataService:  pluginClient,
		metrics:      newMetrics(registerer),
		pCtxProvider: pCtxProvider,
	}
}

// BuildPipeline builds a pipeline from a request.
func (s *Service) BuildPipeline(req *Request) (DataPipeline, error) {
	return s.buildPipeline(req)
}

// ExecutePipeline executes an expression pipeline and returns all the results.
func (s *Service) ExecutePipeline(ctx context.Context, now time.Time, pipeline DataPipeline) (*backend.QueryDataResponse, error) {
	res := backend.NewQueryDataResponse()
	vars, err := pipeline.execute(ctx, now, s)
	if err != nil {
		return nil, err
	}
	for refID, val := range vars {
		res.Responses[refID] = backend.DataResponse{
			Frames: val.Values.AsDataFrames(refID),
			Error:  val.Error,
		}
	}
	return res, nil
}
