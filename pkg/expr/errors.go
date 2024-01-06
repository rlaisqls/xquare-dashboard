package expr

import (
	"errors"
	"fmt"
	"github.com/xquare-dashboard/pkg/services/datasources"
	"github.com/xquare-dashboard/pkg/util/errutil"
)

var unexpectedNodeTypeErrString = "expected executable node type but got node type [{{ .Public.nodeType }} for refid [{{ .Public.refId}}]"

var UnexpectedNodeTypeError = errutil.NewBase(
	errutil.StatusBadRequest, "sse.unexpectedNodeType").MustTemplate(
	unexpectedNodeTypeErrString,
	errutil.WithPublic(unexpectedNodeTypeErrString))

func makeUnexpectedNodeTypeError(refID, nodeType string) error {
	data := errutil.TemplateData{
		Public: map[string]interface{}{
			"refId":    refID,
			"nodeType": nodeType,
		},
		Error: fmt.Errorf("expected executable node type but got node type %v for refId %v", nodeType, refID),
	}

	return UnexpectedNodeTypeError.Build(data)
}

var ConversionError = errutil.BadRequest("sse.readDataError").MustTemplate(
	"[{{ .Public.refId }}] got error: {{ .Error }}",
	errutil.WithPublic(
		"failed to read data from from query {{ .Public.refId }}: {{ .Public.error }}",
	),
)

func makeConversionError(refID string, err error) error {
	data := errutil.TemplateData{
		// Conversion errors should only have meta information in errors
		Public: map[string]any{
			"refId": refID,
			"error": err.Error(),
		},
		Error: err,
	}
	return ConversionError.Build(data)
}

var QueryError = errutil.BadRequest("sse.dataQueryError").MustTemplate(
	"failed to execute query [{{ .Public.refId }}]: {{ .Error }}",
	errutil.WithPublic(
		"failed to execute query [{{ .Public.refId }}]: {{ .Public.error }}",
	))

func MakeQueryError(refID string, datasourceType datasources.DataSourceType, err error) error {
	var pErr error
	var utilErr errutil.Error
	// See if this is grafana error, if so, grab public message
	if errors.As(err, &utilErr) {
		pErr = utilErr.Public()
	} else {
		pErr = err
	}

	data := errutil.TemplateData{
		Public: map[string]any{
			"refId":          refID,
			"datasourceType": datasourceType,
			"error":          pErr.Error(),
		},
		Error: err,
	}

	return QueryError.Build(data)
}
