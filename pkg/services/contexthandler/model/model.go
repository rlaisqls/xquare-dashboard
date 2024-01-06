package contextmodel

import (
	"github.com/xquare-dashboard/pkg/infra/log"
	"github.com/xquare-dashboard/pkg/web"
	"net/http"
)

type ReqContext struct {
	*web.Context
	Logger log.Logger
}

// WriteErr writes an error response based on errutil.Error.
// If provided error is not errutil.Error a 500 response is written.
func (ctx *ReqContext) WriteErr(err error) {
	ctx.writeErrOrFallback(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err)
}

// WriteErrOrFallback uses the information in an errutil.Error if available
// and otherwise falls back to the status and message provided as arguments.
func (ctx *ReqContext) WriteErrOrFallback(status int, message string, err error) {
	ctx.writeErrOrFallback(status, message, err)
}

func (ctx *ReqContext) writeErrOrFallback(status int, message string, err error) {
	data := make(map[string]interface{})
	statusResponse := status

	if err != nil {
		var logMessage string
		logger := ctx.Logger.Warn

		if message != "" {
			logMessage = message
		} else {
			logMessage = http.StatusText(status)
			data["message"] = logMessage
		}

		if status == http.StatusInternalServerError {
			logger = ctx.Logger.Error
		}

		logger(logMessage, "error", err, "remote_addr", ctx.RemoteAddr())
	}

	if _, ok := data["message"]; !ok && message != "" {
		data["message"] = message
	}

	ctx.JSON(statusResponse, data)
}
<<<<<<< HEAD
=======
