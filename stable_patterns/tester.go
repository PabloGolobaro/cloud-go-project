package stable_patterns

import (
	"context"
	"fmt"
	"log"
	"sync"
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
func FanIn_test() {
	sources := make([]<-chan int, 0)
	for i := 0; i < 3; i++ {
		ch := make(chan int)
		sources = append(sources, ch)
		go func() {
			defer close(ch)
			for i := 0; i < 5; i++ {
				ch <- i
				time.Sleep(time.Second)
			}
		}()
	}
	dest := Funnel(sources...)
	for d := range dest {
		fmt.Println(d)

	}
}
func FanOut_test() {
	source := make(chan int)
	dests := Split(source, 5)
	go func() {
		for i := 1; i < 10; i++ {
			source <- i
		}
		close(source)
	}()
	var wg sync.WaitGroup
	wg.Add(len(dests))
	for i, ch := range dests {
		go func(i int, d <-chan int) {
			defer wg.Done()
			for val := range d {
				fmt.Printf("#%d got %d\n", i, val)
			}
		}(i, ch)
	}
	wg.Wait()
}
func ShardedMap_test() {
	shardedMap := NewShardedMap(5)
	shardedMap.Set("alpha", 1)
	shardedMap.Set("beta", 2)
	shardedMap.Set("gama", 3)
	shardedMap.Set("teta", 4)
	fmt.Println(shardedMap.Get("alpha"))
	fmt.Println(shardedMap.Get("beta"))
	fmt.Println(shardedMap.Get("gama"))
	fmt.Println(shardedMap.Get("teta"))
	keys := shardedMap.Keys()
	for _, key := range keys {
		fmt.Println(key)
	}
}
