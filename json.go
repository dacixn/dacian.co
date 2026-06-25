package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type LinkPage struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     []Link `json:"content"`
}

type Link struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

func Parse() (map[string]LinkPage, error) {
	var pages []LinkPage
	file, err := os.Open("static/link-pages.json")
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	dec := json.NewDecoder(file)
	err = dec.Decode(&pages)
	if err != nil {
		return nil, fmt.Errorf("could not decode json: %w", err)
	}
	return sliceToMapByName(pages), nil
}

func sliceToMapByName(pages []LinkPage) map[string]LinkPage {
	var pagesMap map[string]LinkPage
	for _, v := range pages {
		pagesMap[v.Name] = v
	}
	return pagesMap
}
