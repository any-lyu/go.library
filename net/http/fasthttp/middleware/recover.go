package middleware

import (
	"github.com/valyala/fasthttp"

	"github.com/any-lyu/go.library/logs"
	"github.com/any-lyu/go.library/runtime"
)

// Recover 异常处理
func Recover(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				//buf := make([]byte, 64<<10)
				//buf = buf[:runtime.Stack(buf, false)]
				stackStr := runtime.StackString(runtime.Callers(4))
				logs.Error("errgroup: panic recovered: %v\n%s", r, stackStr)
				ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
				ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
				return
			}
		}()
		h(ctx)
	}
}
