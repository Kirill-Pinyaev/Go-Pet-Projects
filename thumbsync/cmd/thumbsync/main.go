package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"thumbsync/internal/jobqueue"
	"thumbsync/internal/worker"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func main() {
	inPath := flag.String("input", "urls.txt", "path to file containing image URLs, one per line")
	outDir := flag.String("out", "./thumbs", "output directory for thumbnails")
	workers := flag.Int("workers", 4, "number of parallel download workers")
	flag.Parse()

	urls, err := readLines(*inPath)
	if err != nil {
		log.Fatalf("failed to read %s: %v", *inPath, err)
	}
	if len(urls) == 0 {
		log.Fatalf("%s is empty", *inPath)
	}

	q := jobqueue.New()
	//pool := worker.NewPool(q, *outDir)

	// enqueue jobs
	go func() {
		for _, u := range urls {
			q.Push(u)
		}
		q.Close()
	}()

	//progress bar
	p := mpb.New(mpb.WithWidth(40))
	bar := p.New(
		int64(len(urls)),
		mpb.BarStyle().Rbound("|"),
		mpb.PrependDecorators(
			decor.Name("completed "),
			decor.CountersNoUnit("%d / %d"),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)

	pool := worker.NewPool(q, *outDir, bar) // ← передаём bar

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool.Run(*workers)

	// ждём, пока либо всё скачалось, либо Ctrl‑C
	done := make(chan struct{})
	go func() { pool.Wait(); close(done) }()

	select {
	case <-ctx.Done():
		q.Close()
		pool.Wait()
	case <-done:
	}

	p.Wait() // прогресс‑бар точно закроется
	fmt.Printf("\nDownloaded: %d, Failed: %d\n", pool.Ok, pool.Fail)
	if len(pool.Errors()) > 0 {
		errFile := filepath.Join(*outDir, "errors.log")
		dumpErrors(errFile, pool.Errors())
		fmt.Printf("Errors saved to %s\n", errFile)
	}
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}
	return lines, s.Err()
}

func dumpErrors(path string, m map[string]error) {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	f, err := os.Create(path)
	if err != nil {
		log.Printf("cannot create error log: %v", err)
		return
	}
	defer f.Close()
	for u, e := range m {
		fmt.Fprintf(f, "%s\t%s\n", u, e)
	}
}
