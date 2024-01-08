package api

import (
	contextmodel "github.com/xquare-dashboard/pkg/services/contexthandler/model"
	"net/http"
)

func (hs *HTTPServer) NotFoundHandler(c *contextmodel.ReqContext) {

	c.JsonApiErr(http.StatusNotFound, "Not found", nil)
}
