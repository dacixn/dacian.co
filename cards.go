package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Card struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type CardPage struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Cards       []Card `json:"cards"`
}

func getCardPagesFromJson(src string) ([]CardPage, error) {
	var pages []CardPage
	file, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	d := json.NewDecoder(file)
	err = d.Decode(&pages)
	if err != nil {
		return nil, fmt.Errorf("could not decode json: %w", err)
	}
	return pages, nil
}

func getCardPageByName(pages []CardPage, name string) (CardPage, error) {
	for _, v := range pages {
		if v.Name == strings.ToLower(name) {
			return v
		}
	}
	return CardPage{}, fmt.Errorf("failed to locate page '%s'", name)
}
