package worker

import (
	"sync"
)

func worker(id, from, to int) int {
	cnt := 0
	for n := from; n <= to; n++ {
		if isPrime(n) {
			cnt++
		}
	}
	return cnt
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func Dispatcher(max, workers int) int {
	results := make(chan int, workers)
	var wg sync.WaitGroup
	wg.Add(workers)
	count := max / workers
	for i := 0; i < workers; i++ {
		start := count * i
		end := count * (i + 1)
		if i == workers-1 {
			end = max
		}
		go func(id, f, t int) {
			defer wg.Done()
			results <- worker(id, f, t)
		}(i, start, end)

	}
	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	for i := range results {
		total += i
	}
	return total
}
