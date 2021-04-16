package middleware

import (
	"net/http"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/any-lyu/go.library/logs"
)

// BaseHandler  log + cross
func BaseHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "content-type")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "OPTIONS,HEAD,GET,POST,PUT,DELETE")
		if string(ctx.Request.Header.Method()) == fasthttp.MethodOptions {
			ctx.SetStatusCode(http.StatusNoContent)
			return
		}
		defer func(t time.Time) {
			since := time.Since(t)
			var body = ""
			if string(ctx.Request.Header.Method()) == fasthttp.MethodPost ||
				string(ctx.Request.Header.Method()) == fasthttp.MethodPut {
				body = string(ctx.Request.Body())
			}
			var method = string(ctx.Request.Header.Method())
			var statusCode = ctx.Response.StatusCode()
			ctx.Response.StatusCode()
			logs.Info("[%s %s %s] %s%d%s [%v] %s %s",
				colorForMethod(method), method, reset,
				colorForStatus(statusCode), statusCode, reset,
				since,
				ctx.Request.RequestURI(),
				body)
		}(time.Now())
		h(ctx)
	}
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}
func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}
