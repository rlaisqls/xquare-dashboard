package datasources

import "github.com/xquare-dashboard/pkg/services/query"

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
	URL:  "",
}

var Prometheus = DataSource{
	Type: LokiType,
	URL:  "",
}

func GetDataSource(dsType DataSourceType) (*DataSource, error) {
	if dsType == LokiType {
		return &Loki, nil
	} else if dsType == PrometheusType {
		return &Prometheus, nil
	} else {
		return nil, query.ErrInvalidDatasourceID
	}
}
