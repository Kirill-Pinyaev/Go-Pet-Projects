package main

import (
	"flag"
	"fmt"
	"microlog/logger"
)

func main() {
	driver := flag.String("driver", "console", "console|buffer")
	dumpAfter := flag.Int("dump-after", 0, "dump buffer after N logs (buffer driver only)")
	anyDemo := flag.Bool("log-any", false, "вывести примеры LogAny")
	flag.Parse()

	var log logger.Logger
	switch *driver {
	case "buffer":
		log = &logger.BufferLogger{}
	default:
		log = &logger.ConsoleLogger{}
	}
	if *anyDemo {
		logger.LogAny("пример строки")
		logger.LogAny(42)
		logger.LogAny(true) // unsupported type bool
	}
	for i := 0; i < 10; i++ {
		log.Info(fmt.Sprintf("msg %d", i))
		if *dumpAfter > 0 && i+1 == *dumpAfter {
			// безопасный вызов Dump() только если есть метод
			if bl, ok := log.(*logger.BufferLogger); ok {
				fmt.Print(bl.Dump())
			}
		}
	}
}
