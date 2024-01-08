package datasources

import (
	"github.com/xquare-dashboard/pkg/util/errutil"
)

type DataSourceType string
type DataSource struct {
	Type DataSourceType
	URL  string
}

const (
	LokiType       DataSourceType = "LokiType"
	PrometheusType DataSourceType = "prometheus"
)

var Loki = DataSource{
	Type: LokiType,
	URL:  "http://loki.xquare.app",
}

var Prometheus = DataSource{
	Type: PrometheusType,
	URL:  "http://prometheus.xquare.app",
}

func GetDataSource(dsType DataSourceType) (*DataSource, error) {
	if dsType == LokiType {
		return &Loki, nil
	} else if dsType == PrometheusType {
		return &Prometheus, nil
	} else {
		return nil, ErrInvalidDatasourceID
	}
}

var (
	ErrInvalidDatasourceID = errutil.BadRequest("query.invalidDatasourceId", errutil.WithPublicMessage("Query does not contain a valid data source identifier")).Errorf("invalid data source identifier")
)
