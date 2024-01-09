package datasources

import (
	"github.com/xquare-dashboard/pkg/util/errutil"
	"os"
)

type DataSourceType string
type DataSource struct {
	Type DataSourceType
	URL  string
}

const (
	LokiType       DataSourceType = "loki"
	PrometheusType DataSourceType = "prometheus"
)

var Loki = DataSource{
	Type: LokiType,
	URL:  os.Getenv("LOKI_URL"),
}

var Prometheus = DataSource{
	Type: PrometheusType,
	URL:  os.Getenv("PROMETHEUS_URL"),
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
