package middleware

import (
	"github.com/any-lyu/go.library/jwt"
	"github.com/valyala/fasthttp"
)

// BasicAuth is the basic auth handler
func BasicAuth(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Get the Basic Authentication credentials
		tokenBytes := ctx.Request.Header.Peek("Authorization")
		token, err := jwt.TokenParse(string(tokenBytes))
		if err != nil {
			// Request Basic Authentication otherwise
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
			ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
			return
		}
		ctx.SetUserValue("uid", token.UID)
		h(ctx)
	}
}
