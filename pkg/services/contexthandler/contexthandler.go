package contexthandler

import (
	"context"
	"github.com/xquare-dashboard/pkg/api/response"
	"github.com/xquare-dashboard/pkg/services/contexthandler/ctxkey"
	contextmodel "github.com/xquare-dashboard/pkg/services/contexthandler/model"
	"github.com/xquare-dashboard/pkg/web"
	"net/http"
)

func ProvideService() *ContextHandler {
	return &ContextHandler{}
}

type ContextHandler struct{}

type reqContextKey = ctxkey.Key

// FromContext returns the ReqContext value stored in a context.Context, if any.
func FromContext(c context.Context) *contextmodel.ReqContext {
	if reqCtx, ok := c.Value(reqContextKey{}).(*contextmodel.ReqContext); ok {
		return reqCtx
	}
	return nil
}

// CopyWithReqContext returns a copy of the parent context with a semi-shallow copy of the ReqContext as a value.
// The ReqContexts's *web.Context is deep copied so that headers are thread-safe; additional properties are shallow copied and should be treated as read-only.
func CopyWithReqContext(ctx context.Context) context.Context {
	origReqCtx := FromContext(ctx)
	if origReqCtx == nil {
		return ctx
	}

	webCtx := &web.Context{
		Req:  origReqCtx.Req.Clone(ctx),
		Resp: web.NewResponseWriter(origReqCtx.Req.Method, response.CreateNormalResponse(http.Header{}, []byte{}, 0)),
	}
	reqCtx := &contextmodel.ReqContext{
		Context: webCtx,
		Logger:  origReqCtx.Logger,
	}
	return context.WithValue(ctx, reqContextKey{}, reqCtx)
}
