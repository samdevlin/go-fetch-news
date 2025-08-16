package main

import (
	"encoding/csv"
	"errors"
	"os"
)

type source struct {
	name string
	url  string
}

func main() {
	loadCSV("sources.csv")
}

func loadCSV(filename string) ([]source, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("the specified file could not be opened")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("the specified file could not be parsed")
	}

	var res []source

	for _, value := range records {
		res = append(res, source{value[0], value[1]})
	}

	return res, nil
}
