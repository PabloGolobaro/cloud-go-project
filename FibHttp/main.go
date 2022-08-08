package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
)

const (
	jaegerEndpoint = "http://localhost:14268/api/traces"
	serviceName    = "fibonacci"
)

func init() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = "" // Отключить вывод отметок времени
	cfg.Sampling = &zap.SamplingConfig{
		Initial:    3, // Допускается регистрировать до трех событий в секунду
		Thereafter: 3, // после превышения предела регистрировать
		// только 1 событие из 3

		Hook: func(e zapcore.Entry, d zapcore.SamplingDecision) {
			if d == zapcore.LogDropped {
				fmt.Println("event dropped...")
			}
		},
	}
	logger, _ := cfg.Build()   // Создать новый регистратор
	zap.ReplaceGlobals(logger) // Заменить глобальный регистратор Zap
}
func main() {

	////////////////////////////////////////////////////
	config := prometheus.Config{}
	ctrl := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)

	Prometheusexporter, err := prometheus.New(config, ctrl)
	if err != nil {
		panic(err)
	}
	mp := Prometheusexporter.MeterProvider()
	global.SetMeterProvider(mp)
	////////////////////////////////////////////
	err = buildRequestsCounter()
	if err != nil {
		log.Fatal(err)
	}
	buildruntimeobservers(context.Background())
	////////////////////////////////////////////////////
	// Создать и настроить консольный экспортер
	_, err = stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps())
	if err != nil {
		log.Fatal(err)
	}
	// Создать и настроить экспортер для Jaeger
	_, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		log.Fatal(err)
	}
	// Создать провайдера трассировки и зарегистрировать в нем
	// вновь созданных экспортеров.
	tp := trace.NewTracerProvider(
	//trace.WithSyncer(stdExporter),
	//trace.WithSyncer(jaegerExporter)
	)
	// Теперь можно зарегистрировать tp как провайдер трассировки otel.
	otel.SetTracerProvider(tp)
	// Зарегистрировать и инструментировать обработчик службы
	http.Handle("/",
		otelhttp.NewHandler(http.HandlerFunc(fibHandler), "root"))
	http.Handle("/metrics", Prometheusexporter)
	// Запустить службу на порту 3000
	log.Fatal(http.ListenAndServe(":3000", nil))
}
