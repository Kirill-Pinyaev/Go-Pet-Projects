package main

import (
	"csvgrep/filter"
	"csvgrep/reader"
	"fmt"
	"os"
)

func main() {
	data := reader.ReadFile(os.Args[len(os.Args)-1])
	err := filter.Filter(data)
	if err != nil {
		fmt.Println(err)
	}
}
