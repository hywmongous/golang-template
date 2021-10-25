package jaeger

import (
	"context"
	"io"

	"github.com/cockroachdb/errors"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

const (
	// https://github.com/opentracing/specification/blob/master/semantic_conventions.mdevent
	ErrorTag       = "error"
	ErrorReportTag = "error.report"

	jaegerUDPPort = 32000
)

var (
	ErrFailedOpeningTransport = errors.New("opening transport for jaeger failed")
	ErrSpanNotFoundInContext  = errors.New("span was not found in context")
	ErrTracerNotFoundInSpan   = errors.New("tracer was not found in span")
)

func Create() (opentracing.Tracer, io.Closer) {
	serviceName := "golang-template"
	jaegerHostPort := "jaeger:6831"

	sampler := jaeger.NewConstSampler(true)

	transport, err := jaeger.NewUDPTransport(jaegerHostPort, jaegerUDPPort)
	if err != nil {
		panic(errors.Wrap(err, ErrFailedOpeningTransport.Error()))
	}

	reporter := jaeger.NewRemoteReporter(transport)
	propagator := zipkin.NewZipkinB3HTTPHeaderPropagator()

	tracer, closer := jaeger.NewTracer(
		serviceName,
		sampler,
		reporter,
		jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, propagator),
		jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, propagator),
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
	)

	if !opentracing.IsGlobalTracerRegistered() {
		opentracing.SetGlobalTracer(tracer)
	}

	return tracer, closer
}

func StartSpanFromSpanContext(
	ctx context.Context,
	operationName string,
	opts ...opentracing.StartSpanOption,
) (opentracing.Span, context.Context) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		panic(errors.Wrap(ErrSpanNotFoundInContext, operationName))
	}

	tracer := span.Tracer()
	if tracer == nil {
		panic(errors.Wrap(ErrTracerNotFoundInSpan, operationName))
	}

	return opentracing.StartSpanFromContextWithTracer(ctx, tracer, operationName)
}

func SetError(span opentracing.Span, err error) {
	span.SetTag(ErrorTag, true)
	span.SetTag(ErrorReportTag, err)
}
