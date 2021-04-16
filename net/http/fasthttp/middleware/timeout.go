package middleware

import (
	"time"

	"github.com/valyala/fasthttp"
)

//Timeout http timeout
func Timeout(h fasthttp.RequestHandler, timeout time.Duration, msg string) fasthttp.RequestHandler {
	return fasthttp.TimeoutHandler(h, timeout, msg)
}
