package csv

import (
	"encoding/csv"
	"io"
	"os"
)

func ReadFile(path string, c chan Line, separator rune) {
	index := 0
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		c <- Line{Error: err}
		close(c)
		return
	}

	reader := csv.NewReader(file)
	reader.Comma = separator
	reader.LazyQuotes = true

	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				c <- Line{Error: err}
				close(c)
				return
			}
		} else {
			if index > 0 && len(line) > 0 {
				c <- Line{
					Index:    index,
					Elements: line,
				}
			}
			index++
		}
	}

	close(c)
}
