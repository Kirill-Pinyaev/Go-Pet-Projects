package main

import (
	"bufio"
	"crawler-chan/internal/fetch"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	var (
		workers = flag.Int("workers", 4, "goroutines extracting links")
		bufsize = flag.Int("bufsize", 100, "capacity of jobs channel")
		input   = flag.String("input", "", "file with seed URLs (stdin if empty)")
	)
	flag.Parse()

	// каналы
	jobs := make(chan string, *bufsize)
	results := make(chan fetch.Result)
	done := make(chan struct{})
	visited := &fetch.Visited{}

	// wait group для воркеров
	var wg sync.WaitGroup
	wg.Add(*workers)
	for i := 0; i < *workers; i++ {
		go fetch.Worker(i, jobs, results, &wg, done, visited)
	}

	// producer
	go func() {
		r := os.Stdin
		if *input != "" {
			f, err := os.Open(*input)
			if err != nil {
				log.Fatalf("open input: %v", err)
			}
			defer f.Close()
			r = f
		}
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			line := sc.Text()
			if line == "" {
				continue
			}
			jobs <- line
		}
		if err := sc.Err(); err != nil {
			log.Printf("scan: %v", err)
		}
	}()

	// агрегируем статистику
	var pages, links, errs int
	timeout := 2 * time.Second
	timer := time.NewTimer(timeout)

loop:
	for {
		select {
		case res := <-results:
			if res.URL == "" && res.Err == nil {
				// канал закрыт
				break loop
			}
			if res.Err != nil {
				errs++
			}
			pages++
			links += res.Links
			log.Printf("%s [%d] %d байт (links %d) err=%v",
				res.URL, res.Status, res.Length, res.Links, res.Err)
			timer.Reset(timeout)

		case <-timer.C:
			log.Print("no activity 2 s — graceful stop")
			close(done) // ❸ сигнал воркерам «стоп отправлять»
			close(jobs) // ❹ теперь безопасно закрыть jobs
			break loop
		}
	}

	// ждём завершения воркеров и закрываем results
	wg.Wait()
	close(results)

	fmt.Printf("\nParsed: %d pages\nLinks:  %d collected\nErrors: %d\n",
		pages, links, errs)
}
