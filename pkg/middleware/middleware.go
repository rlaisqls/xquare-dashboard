package middleware

import (
	"context"
	"github.com/xquare-dashboard/pkg/web"
	"net/http"
)

type contextKey struct{}

var routeOperationNameKey = contextKey{}

func ProvideRouteOperationName(name string) web.Handler {
	return func(res http.ResponseWriter, req *http.Request, c *web.Context) {
		c.Req = addRouteNameToContext(c.Req, name)
	}
}

func addRouteNameToContext(req *http.Request, operationName string) *http.Request {
	// don't set route name if it's set
	if _, exists := RouteOperationName(req); exists {
		return req
	}

	ctx := context.WithValue(req.Context(), routeOperationNameKey, operationName)
	return req.WithContext(ctx)
}

// RouteOperationName receives the route operation name from context, if set.
func RouteOperationName(req *http.Request) (string, bool) {
	if val := req.Context().Value(routeOperationNameKey); val != nil {
		op, ok := val.(string)
		return op, ok
	}

	return "", false
}
