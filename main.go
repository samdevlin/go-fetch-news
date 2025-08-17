package main

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type source struct {
	name string
	url  string
}

type Urlset struct {
	Url []Url `xml:"url"`
}

type Url struct {
	Loc  string `xml:"loc"`
	News News   `xml:"news:news"`
}

type News struct {
	Title    string `xml:"news:title"`
	Keywords string `xml:"news:keywords"`
}

func main() {
	sources, err := loadCSV("sources.csv")
	if err != nil {
		fmt.Println("A fatal error occurred. Exiting..")
	}

	for _, s := range sources {
		// Fetch the xml
		resp, err := http.Get(s.url)
		if err != nil {
			fmt.Printf("An error occurred fetching source %s, skipping..", s.name)
			continue
		}

		// Parse & process
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var urlset Urlset
		if err := xml.Unmarshal(data, &urlset); err != nil {
			panic(err)
		}

		// TODO: emit an event
		for _, item := range urlset.Url {
			fmt.Printf("%s ==> %s\n", item.News.Title, item.Loc)
		}
	}
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
		// TODO: handle out of range
		res = append(res, source{value[0], value[1]})
	}

	return res, nil
}
