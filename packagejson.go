package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Dependencies map[string]string

type PackageJSON struct {
	Name                 string       `json:"name"`
	Version              string       `json:"version"`
	Dependencies         Dependencies `json:"dependencies"`
	DevDependencies      Dependencies `json:"devDependencies"`
	PeerDependencies     Dependencies `json:"peerDependencies"`
	BundledDependencies  Dependencies `json:"bundledDependencies"`
	OptionalDependencies Dependencies `json:"optionalDependencies"`
}

func ReadPackageJSON(filename string) (*PackageJSON, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %w", filename, err)
	}

	var packageJSON PackageJSON
	err = json.Unmarshal(data, &packageJSON)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling data from %q: %w", filename, err)
	}

	return &packageJSON, nil
}
