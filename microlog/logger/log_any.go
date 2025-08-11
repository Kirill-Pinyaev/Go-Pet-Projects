package logger

import "fmt"

func LogAny(v interface{}) {
	switch val := v.(type) {
	case string:
		fmt.Println("string:", val)
	case int, int32, int64:
		fmt.Println("integer:", val)
	default:
		fmt.Printf("unsupported type %T\n", val)
	}
}
