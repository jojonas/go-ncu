package main

import (
	"io/ioutil"
	"testing"
)

func TestLoadResponse(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/api-express.json")
	if err != nil {
		t.Fatalf("reading testdata: %v", err)
	}

	pkg, err := ParsePackageResponse(data)
	if err != nil {
		t.Fatalf("parsing testdata: %v", err)
	}

	if pkg.Name != "express" {
		t.Errorf("expected name 'express', got %q", pkg.Name)
	}
}

func TestAPI(t *testing.T) {
	pkg, err := GetPackage("express")
	if err != nil {
		t.Fatalf("getting 'express' from API: %v", err)
	}
	if pkg.Name != "express" {
		t.Errorf("expected name 'express', got %q", pkg.Name)
	}
}
