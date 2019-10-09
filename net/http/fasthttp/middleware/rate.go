package middleware

import (
	"github.com/any-lyu/go.library/errors"
	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
	"time"
)

// RateHandler 限流
// bursts of at most b tokens.
func RateHandler(h fasthttp.RequestHandler, b int) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		limiter := rate.NewLimiter(rate.Every(time.Second), b)
		if limiter.Allow() {
			h(ctx)
			return
		}
		ctx.Error(errors.ErrSystemBusy.Error(), errors.ErrCode(errors.ErrSystemBusy))
		return
	}
}
