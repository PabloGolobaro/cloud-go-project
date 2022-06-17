package stable_patterns

import (
	"context"
	"fmt"
	"log"
	"time"
)

func Retry_test() {
	fmt.Println("Retry_Effector")
	retry := Retry(EmulateTransienrError, 5, 2*time.Second)
	result, err := retry(context.Background())
	fmt.Println(result, err)

}
func Debounce_Breaker_test() {
	fmt.Println("Debounce_Breaker")
	wrapped := Breaker(DebounceFirst(myFunction, time.Second*1), 3)
	responce, err := wrapped(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println(responce)
}
