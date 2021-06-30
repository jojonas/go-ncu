package main

import (
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestCompare(t *testing.T) {
	type TestCase struct {
		From     string
		To       string
		Expected int
	}

	var tests []TestCase = []TestCase{
		{"1.2.3", "1.2.3", 0},
		{"1.2.3", "2.0.0", 4},
		{"1.2.5", "1.3.0", 3},
		{"1.0.0", "1.0.1", 2},
		{"1.0.0-beta.1", "1.0.0-beta.2", 1},
	}

	for _, test := range tests {
		result := Compare(semver.MustParse(test.From), semver.MustParse(test.To))
		if result != test.Expected {
			t.Errorf("comparing %s -> %s, expected %d, got %d", test.From, test.To, test.Expected, result)
		}

		result = Compare(semver.MustParse(test.To), semver.MustParse(test.From))
		expected := -test.Expected // inverted
		if result != expected {
			t.Errorf("comparing %s -> %s, expected %d, got %d", test.To, test.From, expected, result)
		}
	}
}
