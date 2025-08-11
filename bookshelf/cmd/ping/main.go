package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	addr := flag.String("addr", "http://localhost:8080", "API address")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, *addr+"/ping", nil)

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	d := time.Since(start)

	if err != nil {
		fmt.Printf("request failed after %s: %v\n", d, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("/ping %d in %s\n", resp.StatusCode, d)
}
