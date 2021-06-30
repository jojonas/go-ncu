package main

import (
	"testing"
)

func TestReadLarge(t *testing.T) {
	_, err := ReadPackageJSON("testdata/package-large.json")
	if err != nil {
		t.Fatalf("reading large package.json: %v", err)
	}
}
