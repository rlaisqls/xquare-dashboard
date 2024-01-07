package dtos

import (
	"github.com/xquare-dashboard/pkg/components/simplejson"
	"github.com/xquare-dashboard/pkg/infra/log"
	"regexp"
)

var regNonAlphaNumeric = regexp.MustCompile("[^a-zA-Z0-9]+")
var mlog = log.New("models")

type AnyId struct {
	Id int64 `json:"id"`
}

type AnalyticsSettings struct {
	Identifier         string `json:"identifier"`
	IntercomIdentifier string `json:"intercomIdentifier,omitempty"`
}

type UserPermissionsMap map[string]bool

// swagger:model
type MetricRequest struct {
	// From Start time in epoch timestamps in milliseconds or relative using Grafana time units.
	// required: true
	// example: now-1h
	From string `json:"from"`
	// To End time in epoch timestamps in milliseconds or relative using Grafana time units.
	// required: true
	// example: now
	To string `json:"to"`
	// queries.refId – Specifies an identifier of the query. Is optional and default to “A”.
	// queries.datasourceId – Specifies the data source to be queried. Each query in the request must have an unique datasourceId.
	// queries.maxDataPoints - Species maximum amount of data points that dashboard panel can render. Is optional and default to 100.
	// queries.intervalMs - Specifies the time interval in milliseconds of time series. Is optional and defaults to 1000.
	// required: true
	// example: [ { "refId": "A", "intervalMs": 86400000, "maxDataPoints": 1092, "datasource":{ "uid":"PD8C576611E62080A" }, "rawSql": "SELECT 1 as valueOne, 2 as valueTwo", "format": "table" } ]
	Queries []*simplejson.Json `json:"queries"`
	// required: false
	Debug bool `json:"debug"`
}

func (mr *MetricRequest) GetUniqueDatasourceTypes() []string {
	dsTypes := make(map[string]bool)
	for _, query := range mr.Queries {
		if dsType, ok := query.Get("datasource").CheckGet("type"); ok {
			name := dsType.MustString()
			if _, ok := dsTypes[name]; !ok {
				dsTypes[name] = true
			}
		}
	}

	res := make([]string, 0, len(dsTypes))
	for dsType := range dsTypes {
		res = append(res, dsType)
	}

	return res
}

func (mr *MetricRequest) CloneWithQueries(queries []*simplejson.Json) MetricRequest {
	return MetricRequest{
		From:    mr.From,
		To:      mr.To,
		Queries: queries,
		Debug:   mr.Debug,
	}
}
