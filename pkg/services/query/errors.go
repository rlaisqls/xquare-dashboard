package query

import (
	"github.com/xquare-dashboard/pkg/util/errutil"
)

var (
	ErrNoQueriesFound        = errutil.BadRequest("query.noQueries", errutil.WithPublicMessage("No queries found")).Errorf("no queries found")
	ErrMissingDataSourceInfo = errutil.BadRequest("query.missingDataSourceInfo").MustTemplate("query missing datasources info: {{ .Public.RefId }}", errutil.WithPublic("Query {{ .Public.RefId }} is missing datasources information"))
	ErrQueryParamMismatch    = errutil.BadRequest("query.headerMismatch", errutil.WithPublicMessage("The request headers point to a different plugin than is defined in the request body")).Errorf("plugin header/body mismatch")
	ErrDuplicateRefId        = errutil.BadRequest("query.duplicateRefId", errutil.WithPublicMessage("Multiple queries using the same RefId is not allowed ")).Errorf("multiple queries using the same RefId is not allowed")
)
