package stable_patterns

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func Throttle(effector Effector, max uint, refill uint, duration time.Duration) Effector {
	var tokens = max
	var once sync.Once

	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		once.Do(func() {
			ticker := time.NewTicker(duration)
			go func() {
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
						log.Println("Token is added: ", tokens)

					}
				}
			}()
		})
		if tokens <= 0 {
			return "", fmt.Errorf("to many calls")
		}
		tokens--
		log.Println("Token is been used: ", tokens)
		return effector(ctx)
	}

}
