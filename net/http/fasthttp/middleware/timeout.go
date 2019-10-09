package middleware

import (
	"github.com/valyala/fasthttp"
	"time"
)

//Timeout http timeout
func Timeout(h fasthttp.RequestHandler, timeout time.Duration, msg string) fasthttp.RequestHandler {
	return fasthttp.TimeoutHandler(h, timeout, msg)
}
