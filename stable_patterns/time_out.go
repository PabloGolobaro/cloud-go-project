package stable_patterns

import (
	"context"
	"time"
)

type SlowFunction func(str string) (string, error)

type WithContext func(ctx context.Context, str string) (string, error)

func TimeOut(function SlowFunction) WithContext {
	return func(ctx context.Context, str string) (string, error) {
		chres := make(chan string)
		cherr := make(chan error)
		go func() {
			res, err := function(str)
			chres <- res
			cherr <- err
		}()

		select {
		case res := <-chres:
			return res, <-cherr
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

func Slow(string2 string) (string, error) {
	<-time.NewTimer(time.Second * 1).C
	return string2, nil
}
