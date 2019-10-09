package tracing

import (
	"context"
	"github.com/any-lyu/go.library/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net"
	"net/http"
	"strconv"
)

// ContextToHTTP propagate tracing info into HTTP header, useful in server side.
func ContextToHTTP(ctx context.Context, req *http.Request) {
	// Try to find a Span in the Context.
	if span := opentracing.SpanFromContext(ctx); span != nil {
		// Add standard OpenTracing tags.
		ext.HTTPMethod.Set(span, req.Method)
		ext.HTTPUrl.Set(span, req.URL.String())
		host, portString, err := net.SplitHostPort(req.URL.Host)
		if err == nil {
			ext.PeerHostname.Set(span, host)
			if port, err := strconv.Atoi(portString); err != nil {
				ext.PeerPort.Set(span, uint16(port))
			}
		} else {
			ext.PeerHostname.Set(span, req.URL.Host)
		}

		if err = opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.TextMap,
			opentracing.HTTPHeadersCarrier(req.Header),
		); err != nil {
			logs.Error(err.Error())
		}
	} else {
		logs.Error("can not extract span from context")
	}
}

// HTTPToContext for client side service get tracing info from HTTP header.
func HTTPToContext(ctx context.Context, req *http.Request,
	tracer opentracing.Tracer, operationName string) context.Context {
	// Try to join to a trace propagated in `req`.
	var span opentracing.Span
	wireContext, err := tracer.Extract(
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	if err != nil && err != opentracing.ErrSpanContextNotFound {
		logs.Error(err.Error())
	}

	span = tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
	ext.HTTPMethod.Set(span, req.Method)
	ext.HTTPUrl.Set(span, req.URL.String())

	return opentracing.ContextWithSpan(ctx, span)
}
