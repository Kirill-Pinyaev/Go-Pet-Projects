package main

import (
	"flag"
	"fmt"
	"os"
	"primes-runner/worker"
	"runtime"
	"time"
)

func main() {
	start := time.Now()
	max := flag.Int("max", 1, "number")
	workers := flag.Int("workers", 4, "how many workers")
	gomax := flag.Int("gomax", 2, "how many cpu")

	flag.Parse()
	if *workers <= 0 {
		fmt.Println("Error: need workers >= 1")
		os.Exit(1)
	}

	runtime.GOMAXPROCS(*gomax)

	total := worker.Dispatcher(*max, *workers)
	elapsed := time.Since(start)
	fmt.Printf("Count simple number: %v, time: %v\n", total, elapsed)
}
