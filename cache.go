package main

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/patrickmn/go-cache"
)

func ItemsFromFile(filename string) (map[string]cache.Item, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening %q: %w", filename, err)
	}
	defer file.Close()

	var items map[string]cache.Item

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&items)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return items, nil
}

func ItemsToFile(filename string, items map[string]cache.Item) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating %q: %w", filename, err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(items)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
