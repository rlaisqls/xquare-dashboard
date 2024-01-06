// Package contexthandler contains the ContextHandler service.
package contexthandler

import (
	"context"
	"net/http"

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
	return context.WithValue(ctx, reqContextKey{}, reqCtx)
	}
	return context.WithValue(ctx, reqContextKey{}, reqCtx)
}

// Middleware provides a middleware to initialize the request context.
func (h *ContextHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := h.tracer.Start(r.Context(), "Auth - Middleware")
		defer span.End() // this will span to next handlers as well

		reqContext := &contextmodel.ReqContext{
			Context: web.FromContext(ctx), // Extract web context from context (no knowledge of the trace)
			Logger:  log.New("context"),
		}

		// inject ReqContext in the context
		ctx = context.WithValue(ctx, reqContextKey{}, reqContext)

		// Set the context for the http.Request.Context
		// This modifies both r and reqContext.Req since they point to the same value
		*reqContext.Req = *reqContext.Req.WithContext(ctx)

		traceID := tracing.TraceIDFromContext(reqContext.Req.Context(), false)
		if traceID != "" {
			reqContext.Logger = reqContext.Logger.New("traceID", traceID)
		}
		next.ServeHTTP(w, r)
	})
}
>>>>>>> 5ad0219 (init project)
