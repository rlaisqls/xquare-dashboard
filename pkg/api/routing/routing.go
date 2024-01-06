package routing

import (
	"github.com/xquare-dashboard/pkg/api/response"
	contextmodel "github.com/xquare-dashboard/pkg/services/contexthandler/model"
	"github.com/xquare-dashboard/pkg/web"
)

var (
	ServerError = func(err error) response.Response {
		return response.Error(500, "Server error", err)
	}
)

func Wrap(handler func(c *contextmodel.ReqContext) response.Response) web.Handler {
	return func(c *contextmodel.ReqContext) {
		if res := handler(c); res != nil {
			res.WriteTo(c)
		}
	}
}
