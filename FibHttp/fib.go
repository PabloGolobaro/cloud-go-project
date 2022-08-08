package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"os"
)

var labels = []attribute.KeyValue{
	attribute.String("application", serviceName),
	attribute.String("container_id", os.Getenv("HOSTNAME")),
}

func Fibonacci(ctx context.Context, n int) chan int {

	ch := make(chan int)
	go func() {
		tr := otel.GetTracerProvider().Tracer(serviceName)
		cctx, sp := tr.Start(ctx,
			fmt.Sprintf("Fibonacci(%d)", n),
			trace.WithAttributes(attribute.Int("n", n)))
		defer sp.End()
		result := 1
		if n > 1 {
			a := Fibonacci(cctx, n-1)
			b := Fibonacci(cctx, n-2)
			result = <-a + <-b
		}
		sp.SetAttributes(attribute.Int("result", result))
		ch <- result
	}()
	return ch
}
