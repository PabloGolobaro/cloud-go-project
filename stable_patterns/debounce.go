package stable_patterns

import (
	"context"
	"sync"
	"time"
)

func myFunction(ctx context.Context) (string, error) {
	return "responce", nil
}
func DebounceFirst(circuit Circuit, duration time.Duration) Circuit {
	var threshlod time.Time
	var result string
	var err error
	var m sync.Mutex
	return func(ctx context.Context) (string, error) {
		m.Lock()
		defer func() {
			threshlod = time.Now().Add(duration)
			m.Unlock()
		}()
		if time.Now().Before(threshlod) {
			return result, err
		}
		result, err = circuit(ctx)
		return result, err
	}
}
