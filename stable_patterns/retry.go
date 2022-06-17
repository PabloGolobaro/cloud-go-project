package stable_patterns

import (
	"context"
	"errors"
	"log"
	"time"
)

var count int

type Effector func(ctx context.Context) (string, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string, error) {
		for r := 0; ; r++ {
			responce, err := effector(ctx)
			if err == nil || r >= retries {
				return responce, err
			}
			log.Printf("Attempt %d failde; retrying in %v", r+1, delay)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}
func EmulateTransienrError(ctx context.Context) (string, error) {
	count++
	if count <= 3 {
		return "internal fail", errors.New("error")
	} else {
		return "succes", nil
	}
}
