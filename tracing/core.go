package tracing

import (
	"context"
	"github.com/any-lyu/go.library/logs"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracinglog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-client-go/transport"
)

// DefaultSampleRate is the default sample rate
const (
	DefaultSampleRate      = 0.01
	redisOperationName     = "redis"
	mongoOperationName     = "mongo"
	goroutineOperationName = "goroutine"
)

// InitGlobalTracer init a global tracer from a given endpoint
func InitGlobalTracer(serviceName, endpoint string, sampleRate float64) (io.Closer, error) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: sampleRate,
		},
	}

	// Initialize tracer with a logger
	closer, err := cfg.InitGlobalTracer(
		serviceName,
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Gen128Bit(true),
		jaegercfg.Reporter(jaeger.NewRemoteReporter(transport.NewHTTPTransport(endpoint))),
	)

	return closer, err
}

func startSpanFactory(
	ctx context.Context, operationName string, referenceType opentracing.SpanReferenceType) opentracing.Span {

	parentSpan := opentracing.SpanFromContext(ctx)

	var span opentracing.Span
	if parentSpan == nil {
		//log.ErrorContext(ctx, "can not start a span from ctx")
		logs.Error("can not start a span from ctx")
		span = opentracing.GlobalTracer().StartSpan(operationName)
	} else {
		switch referenceType {
		case opentracing.ChildOfRef:
			span = opentracing.GlobalTracer().StartSpan(operationName, opentracing.ChildOf(parentSpan.Context()))
		case opentracing.FollowsFromRef:
			span = opentracing.GlobalTracer().StartSpan(operationName, opentracing.FollowsFrom(parentSpan.Context()))
		}
	}

	return span
}

func finishSpanFactory(span opentracing.Span, err error) {

	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(tracinglog.String("event", "error"), tracinglog.String("message", err.Error()))
	}
	span.Finish()
}

// RedisSpan for span operation in redis
type RedisSpan struct{}

// Redis for external package get RedisSpan struct
func Redis() *RedisSpan {
	return &RedisSpan{}
}

// StartSpan start a span from context before any redis operation
func (RedisSpan) StartSpan(ctx context.Context) opentracing.Span {
	span := startSpanFactory(ctx, redisOperationName, opentracing.ChildOfRef)
	ext.SpanKindRPCClient.Set(span)
	ext.DBType.Set(span, redisOperationName)
	return span
}

// FinishSpan finish a span after redis operation
func (RedisSpan) FinishSpan(span opentracing.Span, err error) {
	finishSpanFactory(span, err)
}

// MongoSpan for span operation in mongo
type MongoSpan struct{}

// Mongo for external package get MongoSpan struct
func Mongo() *MongoSpan {
	return &MongoSpan{}
}

// StartSpan start a span from context before any mongo operation
func (MongoSpan) StartSpan(ctx context.Context) opentracing.Span {
	span := startSpanFactory(ctx, mongoOperationName, opentracing.ChildOfRef)
	ext.SpanKindRPCClient.Set(span)
	ext.DBType.Set(span, mongoOperationName)
	return span
}

// FinishSpan finish a span after mongo operation
func (MongoSpan) FinishSpan(span opentracing.Span, err error) {
	finishSpanFactory(span, err)
}

// GeneralSpan for non db tracing, useful in boot up goroutine.
// default is ChildOfRef reference type
type GeneralSpan struct {
}

// General for external package get GeneralSpan struct
func General() *GeneralSpan {
	return &GeneralSpan{}
}

var (
	defaultOptions = &options{
		referenceType: opentracing.ChildOfRef,
	}
)

// Option for GeneralSpan option
type (
	Option func(span *options)

	options struct {
		referenceType opentracing.SpanReferenceType
	}
)

// WithChildOfSpan set ChildOfRef reference type to the span
// NOTE: the difference of ChildOfRef type and FollowsFromRef type is
// logically at present, same behavior in frontend presention
func WithChildOfSpan() Option {
	return func(span *options) {
		span.referenceType = opentracing.ChildOfRef
	}
}

// WithFollowsFromSpan set FollowsFromRef reference type to the span
// Example usage :
// sp := tracing.General().StartSpan(ctx, tracing.WithFollowsFromSpan())
func WithFollowsFromSpan() Option {
	return func(span *options) {
		span.referenceType = opentracing.FollowsFromRef
	}
}

// StartSpan start a span for general purpose,
// default is ChildOfRef reference type
func (GeneralSpan) StartSpan(ctx context.Context, options ...Option) opentracing.Span {
	opt := defaultOptions
	for _, option := range options {
		option(opt)
	}
	span := startSpanFactory(ctx, goroutineOperationName, opt.referenceType)
	return span
}

// FinishSpan finish a span
func (GeneralSpan) FinishSpan(span opentracing.Span, err error) {
	finishSpanFactory(span, err)
}
