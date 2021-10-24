package jaeger

import (
	"errors"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

var ErrFailedOpeningTransport = errors.New("opening transport for jaeger failed")

func Create() (opentracing.Tracer, io.Closer) {
	serviceName := "golang-template"
	jaegerHostPort := "jaeger:6831"

	sampler := jaeger.NewConstSampler(true)

	transport, err := jaeger.NewUDPTransport(jaegerHostPort, 32000)
	if err != nil {
		panic(err)
	}

	reporter := jaeger.NewRemoteReporter(transport)
	propagator := zipkin.NewZipkinB3HTTPHeaderPropagator()

	return jaeger.NewTracer(
		serviceName,
		sampler,
		reporter,
		jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, propagator),
		jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, propagator),
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
	)
}
