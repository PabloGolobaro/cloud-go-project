package main

import (
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
)

func fibHandler(w http.ResponseWriter, req *http.Request) {

	var err error
	var n int
	if len(req.URL.Query()["n"]) != 1 {
		err = fmt.Errorf("wrong number of arguments")
	} else {
		n, err = strconv.Atoi(req.URL.Query()["n"][0])
	}
	if err != nil {
		http.Error(w, "couldn't parse index n", 400)
		return
	}
	// Получить текущий контекст из входящего запроса
	ctx := req.Context()
	requests.Add(ctx, 1, labels...)
	// Вызвать дочернюю функцию и передать ей контекст запроса.
	result := <-Fibonacci(ctx, n)
	// Получить экземпляр Span, связанный с текущим контекстом, и
	// присоединить параметр и результат в виде атрибутов.
	if sp := trace.SpanFromContext(ctx); sp != nil {
		sp.SetAttributes(
			attribute.Int("parameter", n),
			attribute.Int("result", result))
	}
	// Послать ответ с результатом.
	fmt.Fprintln(w, result)
}
