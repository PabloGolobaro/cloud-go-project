package main

import (
	"context"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"runtime"
)

var requests syncint64.Counter

func buildRequestsCounter() error {
	var err error
	// Получить экземпляр Meter из провайдера метрик.
	meter := global.MeterProvider().Meter(serviceName)
	// Получить инструмент Int64Counter для метрики с именем
	// "fibonacci_requests_total".
	requests, err = meter.SyncInt64().Counter("fibonacci_requests_total",
		instrument.WithDescription("Total number of Fibonacci requests."),
	)
	return err
}
func buildruntimeobservers(ctx context.Context) {
	// Получить экземпляр Meter из провайдера метрик.
	meter := global.MeterProvider().Meter(serviceName)
	// Создать инструмент для получения объема используемой памяти
	// и количества сопрограмм. Значения ошибок игнорируются для краткости.
	m := runtime.MemStats{}
	mem, err := meter.AsyncInt64().UpDownCounter("memory_usage_bytes",
		instrument.WithDescription("Amount of memory used."),
	)
	if err != nil {
		return
	}
	gorutins_num, err := meter.AsyncInt64().UpDownCounter("num_goroutines",
		instrument.WithDescription("Number of running goroutines."),
	)
	if err != nil {
		return
	}
	if err := meter.RegisterCallback([]instrument.Asynchronous{
		mem,
	}, func(ctx context.Context) {
		runtime.ReadMemStats(&m)
		mem.Observe(ctx, int64(m.Sys))
	}); err != nil {
		panic(err)
	}
	if err := meter.RegisterCallback([]instrument.Asynchronous{
		gorutins_num,
	}, func(ctx context.Context) {
		num := runtime.NumGoroutine()
		gorutins_num.Observe(ctx, int64(num))
	}); err != nil {
		panic(err)
	}
}
