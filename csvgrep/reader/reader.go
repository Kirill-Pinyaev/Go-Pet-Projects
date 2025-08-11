package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func ReadFile(fileName string) (data [][]string) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("errror read file %s", err)
	}
	reader := csv.NewReader(f)
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		data = append(data, rec)
	}
	return
}
