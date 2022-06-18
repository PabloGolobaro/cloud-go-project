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
func Throttle_test() {
	fmt.Println("Throttle_Test")
	throttled := Throttle(myFunction, 3, 1, time.Second*8)
	for {
		fmt.Printf("Time of call %v \n", time.Now().Format(time.Stamp))
		result, err := throttled(context.Background())
		if err != nil {
			fmt.Println(err)

		}
		fmt.Println(result)
		<-time.NewTimer(time.Second * 2).C
	}
}
func TimeOut_test() {
	fmt.Println("TimeOut test")
	ctx := context.Background()
	withTimeout, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	timeOut := TimeOut(Slow)
	res, err := timeOut(withTimeout, "success")
	fmt.Println(res, err)

}
