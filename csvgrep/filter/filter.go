package filter

import (
	"flag"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var (
	col = flag.String("col", "", "col name")
	op  = flag.String("op", "eq", "operations eq|ne|gt|lt|match")
	val = flag.String("val", "", "value")
)

func Filter(data [][]string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("compare error: %v", r)
		}
	}()
	flag.Parse()
	if *col == "" || *val == "" {
		flag.Usage()
	}
	colIndex := slices.Index(data[0], *col)
	if colIndex == -1 {
		err = fmt.Errorf("colum not found")
		return
	}

	fmt.Println(strings.Join(data[0], ", "))
	for _, v := range data[1:] {
		if compare(*val, v[colIndex]) {
			fmt.Println(strings.Join(v, ", "))
		}
	}
	return
}

func compare(find string, val string) bool {
	switch *op {
	case "eq":
		return val == find
	case "ne":
		return val != find
	case "gt":
		return val > find
	case "lt":
		return val < find
	case "match":
		re, err := regexp.Compile(find)
		if err != nil {
			panic(fmt.Errorf("bad regexp %q: %w", val, err))
		}
		return re.MatchString(val)
	default:
		panic("unsupported op")
	}
}
