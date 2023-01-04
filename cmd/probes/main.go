package main

import (
	"context"
	"fmt"
	"strconv"

	// "runtime"
	"time"
)

type gauge uint64

func second(ctx context.Context) {
	var cancelled bool = false
	ticker := time.NewTicker(1 * time.Second)
	for !cancelled {

		select {
		case <-ticker.C:
			fmt.Println("Work!")
		case <-ctx.Done():
			{
				fmt.Println("Canceled!")
				// case <-ticker.C:
				cancelled = true
			}
		}

	}
}

func main() {

	//ctx := context.Background()
	//ctx2, cancel := context.WithCancel(ctx)
	//timer := time.NewTimer(12250 * time.Millisecond)
	//go second(ctx2)
	//
	//<-timer.C
	//
	//cancel()
	//time.Sleep(2 * time.Second)
	//fmt.Println("Stopped")
	var v float64 = 1.1415 + 13
	for i := -1; i <= 50; i++ {

		s := strconv.FormatFloat(v, 'E', i, 32)
		//s2 := strconv.FormatFloat(1.1415+13, 'E', i, 64)
		fmt.Println(i, ":", s)

	}

}
