package query

import (
	"context"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/xquare-dashboard/pkg/api/dtos"
	"github.com/xquare-dashboard/pkg/components/simplejson"
	"github.com/xquare-dashboard/pkg/expr"
	"github.com/xquare-dashboard/pkg/infra/log"
	"github.com/xquare-dashboard/pkg/services/contexthandler"
	"github.com/xquare-dashboard/pkg/services/datasources"
	"github.com/xquare-dashboard/pkg/tsdb/legacydata"
	"github.com/xquare-dashboard/pkg/util/errutil"
	"golang.org/x/sync/errgroup"
	"net/http"
	"runtime"
	"slices"
	"time"
)

const (
	HeaderPluginID       = "X-Plugin-Id"          // can be used for routing
	HeaderDatasourceUID  = "X-DataSourceType-Uid" // can be used for routing/ load balancing
	HeaderDashboardUID   = "X-Dashboard-Uid"      // mainly useful for debugging slow queries
	HeaderPanelID        = "X-Panel-Id"           // mainly useful for debugging slow queries
	HeaderQueryGroupID   = "X-Query-Group-Id"     // mainly useful for finding related queries with query chunking
	HeaderFromExpression = "X-Grafana-From-Expr"  // used by datasources to identify expression queries
)

func ProvideService(expressionService *expr.Service) *ServiceImpl {
	g := &ServiceImpl{
		expressionService:    expressionService,
		log:                  log.New("query_data"),
		concurrentQueryLimit: runtime.NumCPU(),
	}
	g.log.Info("Query Service initialization")
	return g
}

type Service interface {
	Run(ctx context.Context) error
	QueryData(ctx context.Context, reqDTO dtos.MetricRequest) (*backend.QueryDataResponse, error)
}

// Gives us compile time error if the service does not adhere to the contract of the interface
var _ Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
	expressionService    *expr.Service
	log                  log.Logger
	concurrentQueryLimit int
}

// Run ServiceImpl.
func (s *ServiceImpl) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

// QueryData processes queries and returns query responses. It handles queries to single or mixed datasources, as well as expressions.
func (s *ServiceImpl) QueryData(ctx context.Context, reqDTO dtos.MetricRequest) (*backend.QueryDataResponse, error) {
	// Parse the request into parsed queries grouped by datasource uid
	parsedReq, err := s.parseMetricRequest(ctx, reqDTO)
	if err != nil {
		return nil, err
	}

	// If there are expressions, handle them and return
	if parsedReq.hasExpression {
		return s.handleExpressions(ctx, parsedReq)
	}

	// If there are multiple datasources, handle their queries concurrently and return the aggregate result
	return s.executeConcurrentQueries(ctx, reqDTO, parsedReq.parsedQueries)
}

// splitResponse contains the results of a concurrent data source query - the response and any headers
type splitResponse struct {
	responses backend.Responses
	header    http.Header
}

// executeConcurrentQueries executes queries to multiple datasources concurrently and returns the aggregate result.
func (s *ServiceImpl) executeConcurrentQueries(
	ctx context.Context, reqDTO dtos.MetricRequest, queriesbyDs map[datasources.DataSourceType][]parsedQuery,
) (*backend.QueryDataResponse, error) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(s.concurrentQueryLimit) // prevent too many concurrent requests
	rchan := make(chan splitResponse, len(queriesbyDs))

	// Create panic recovery function for loop below
	recoveryFn := func(queries []*simplejson.Json) {
		if r := recover(); r != nil {
			var err error
			s.log.Error("query datasource panic", "error", r, "stack", log.Stack(1))
			if theErr, ok := r.(error); ok {
				err = theErr
			} else if theErrString, ok := r.(string); ok {
				err = fmt.Errorf(theErrString)
			} else {
				err = fmt.Errorf("unexpected error")
			}
			// Due to the panic, there is no valid response for any query for this datasource. Append an error for each one.
			rchan <- buildErrorResponses(err, queries)
		}
	}

	// Query each datasource concurrently
	for _, queries := range queriesbyDs {
		rawQueries := make([]*simplejson.Json, len(queries))
		for i := 0; i < len(queries); i++ {
			rawQueries[i] = queries[i].rawQuery
		}
		g.Go(func() error {
			subDTO := reqDTO.CloneWithQueries(rawQueries)
			// Handle panics in the datasource query
			defer recoveryFn(subDTO.Queries)

			ctxCopy := contexthandler.CopyWithReqContext(ctx)
			subResp, err := s.QueryData(ctxCopy, subDTO)
			if err == nil {
				reqCtx, header := contexthandler.FromContext(ctxCopy), http.Header{}
				if reqCtx != nil {
					header = reqCtx.Resp.Header()
				}
				rchan <- splitResponse{subResp.Responses, header}
			} else {
				// If there was an error, return an error response for each query for this datasource
				rchan <- buildErrorResponses(err, subDTO.Queries)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(rchan)
	resp := backend.NewQueryDataResponse()
	reqCtx := contexthandler.FromContext(ctx)
	for result := range rchan {
		for refId, dataResponse := range result.responses {
			resp.Responses[refId] = dataResponse
		}
		if reqCtx != nil {
			for k, v := range result.header {
				for _, val := range v {
					if !slices.Contains(reqCtx.Resp.Header().Values(k), val) {
						reqCtx.Resp.Header().Add(k, val)
					} else {
						s.log.Warn("skipped duplicate response header", "header", k, "value", val)
					}
				}
			}
		}
	}

	return resp, nil
}

// buildErrorResponses applies the provided error to each query response in the list. These queries should all belong to the same datasource.
func buildErrorResponses(err error, queries []*simplejson.Json) splitResponse {
	er := backend.Responses{}
	for _, query := range queries {
		er[query.Get("refId").MustString("A")] = backend.DataResponse{
			Error: err,
		}
	}
	return splitResponse{er, http.Header{}}
}

// handleExpressions handles POST /api/ds/query when there is an expression.
func (s *ServiceImpl) handleExpressions(ctx context.Context, parsedReq *parsedRequest) (*backend.QueryDataResponse, error) {
	exprReq := expr.Request{
		Queries: []expr.Query{},
	}

	for _, pq := range parsedReq.getFlattenedQueries() {
		if pq.datasource == nil {
			return nil, ErrMissingDataSourceInfo.Build(errutil.TemplateData{
				Public: map[string]any{
					"RefId": pq.query.RefID,
				},
			})
		}

		exprReq.Queries = append(exprReq.Queries, expr.Query{
			JSON:          pq.query.JSON,
			Interval:      pq.query.Interval,
			RefID:         pq.query.RefID,
			MaxDataPoints: pq.query.MaxDataPoints,
			QueryType:     pq.query.QueryType,
			DataSource:    pq.datasource,
			TimeRange: expr.AbsoluteTimeRange{
				From: pq.query.TimeRange.From,
				To:   pq.query.TimeRange.To,
			},
		})
	}

	qdr, err := s.expressionService.TransformData(ctx, time.Now(), &exprReq) // use time now because all queries have absolute time range
	if err != nil {
		return nil, fmt.Errorf("expression request error: %w", err)
	}
	return qdr, nil
}

// parseRequest parses a request into parsed queries grouped by datasource uid
func (s *ServiceImpl) parseMetricRequest(ctx context.Context, reqDTO dtos.MetricRequest) (*parsedRequest, error) {
	if len(reqDTO.Queries) == 0 {
		return nil, ErrNoQueriesFound
	}

	timeRange := legacydata.NewDataTimeRange(reqDTO.From, reqDTO.To)
	req := &parsedRequest{
		hasExpression: false,
		parsedQueries: make(map[datasources.DataSourceType][]parsedQuery),
		dsTypes:       make(map[datasources.DataSourceType]bool),
	}

	// Parse the queries and store them by datasource
	datasources := map[datasources.DataSourceType]*datasources.DataSource{}
	for _, query := range reqDTO.Queries {
		ds, err := s.getDataSourceFromQuery(query)
		if err != nil {
			return nil, err
		}

		datasources[ds.Type] = ds
		req.hasExpression = true
		req.dsTypes[ds.Type] = true

		if _, ok := req.parsedQueries[ds.Type]; !ok {
			req.parsedQueries[ds.Type] = []parsedQuery{}
		}

		s.log.Debug("Processing metrics query", "query", query)

		modelJSON, err := query.MarshalJSON()
		if err != nil {
			return nil, err
		}

		req.parsedQueries[ds.Type] = append(req.parsedQueries[ds.Type], parsedQuery{
			datasource: ds,
			query: backend.DataQuery{
				TimeRange: backend.TimeRange{
					From: timeRange.GetFromAsTimeUTC(),
					To:   timeRange.GetToAsTimeUTC(),
				},
				RefID:         query.Get("refId").MustString("A"),
				MaxDataPoints: query.Get("maxDataPoints").MustInt64(100),
				Interval:      time.Duration(query.Get("intervalMs").MustInt64(1000)) * time.Millisecond,
				QueryType:     query.Get("queryType").MustString(""),
				JSON:          modelJSON,
			},
			rawQuery: query,
		})
	}

	return req, req.validateRequest(ctx)
}

func (s *ServiceImpl) getDataSourceFromQuery(query *simplejson.Json) (*datasources.DataSource, error) {
	ds := datasources.DataSourceType(query.Get("datasource").MustString())
	return datasources.GetDataSource(ds)
}
