package main

import (
	"time"

	"github.com/hywmongous/example-service/internal/infrastructure/jaeger"
	"github.com/opentracing/opentracing-go/log"
)

func main() {
	tracer, closer := jaeger.Create()
	defer closer.Close()

	span := tracer.StartSpan("test")
	defer span.Finish()
	span.LogFields(
		log.String("val1", "string"),
	)

	sleepTime := 1 * time.Second
	time.Sleep(sleepTime)
}
