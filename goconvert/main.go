package main

import (
	"fmt"
	"os"
)

func main() {
	prec, num, unit, err := parseArgs()
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
	endNum, endUnit, err := convert(num, unit)
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%.*f %s\n", prec, endNum, endUnit)

}
