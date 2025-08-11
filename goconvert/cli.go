package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)

func parseArgs() (prec int, num float64, unit string, err error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	precPtr := fs.Int("prec", 2, "precision for output")

	if err = fs.Parse(os.Args[1:]); err != nil {
		return
	}
	if *precPtr < 0 {
		err = fmt.Errorf("precision must be non-negative")
		return
	}

	args := fs.Args()
	switch len(args) {
	case 1: // формат "12kg"
		raw := args[0]
		idx := -1
		for i, r := range raw {
			if !unicode.IsDigit(r) && r != '.' && r != '-' && r != '+' &&
				r != 'e' && r != 'E' {
				idx = i
				break
			}
		}
		if idx == -1 {
			err = fmt.Errorf("missing unit in %q", raw)
			return
		}
		if num, err = strconv.ParseFloat(raw[:idx], 64); err != nil {
			err = fmt.Errorf("invalid number %q: %w", raw[:idx], err)
			return
		}
		unit = raw[idx:]

	case 2: // формат "12 kg"
		if num, err = strconv.ParseFloat(args[0], 64); err != nil {
			err = fmt.Errorf("invalid number %q: %w", args[0], err)
			return
		}
		unit = args[1]

	default:
		err = fmt.Errorf("usage: convert <value><unit> # km, mi, kg, lb, C, F [--prec N]")
		return
	}
	prec = *precPtr
	return

}
