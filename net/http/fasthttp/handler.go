package fasthttp

import (
	"github.com/any-lyu/go.library/errors"
	"github.com/any-lyu/go.library/logs"
	"github.com/any-lyu/go.library/tool"
	"github.com/valyala/fasthttp"
)

// HandlerFunc 是 http handler 函数模型.
type HandlerFunc func(c *fasthttp.RequestCtx) (data interface{}, err error)

// HandlerFuncWrapper 返回的统一封装
func HandlerFuncWrapper(fn HandlerFunc) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		type response struct {
			Code    int         `json:"code"`
			Message string      `json:"msg"`
			Data    interface{} `json:"data"`
		}
		resp := response{}
		data, err := fn(ctx)
		if err == nil {
			resp.Data = data
			logs.Debug("Response:", tool.D2S(data))
			goto back
		}
		logs.Debug("ResponseErr:", err.Error())
		err = errors.Cause(err)
		switch err {
		case errors.ErrSystemBusy:
			ctx.Error(errors.ErrSystemBusy.Error(), errors.ErrCode(errors.ErrSystemBusy))
			return
		case errors.ErrToken:
			ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
			return
		default:
			resp.Message, resp.Code = errors.ErrCodeMessage(err)
		}
	back:
		ctx.Success("application/json", tool.D2B(data))
	}
}
