package ynabler

import (
	"testing"
)

func TestCiti(t *testing.T) {
	file := "test-data/citi.csv"

	lines, err := ReadFile(file)
	if err != nil {
		panic(err)
	}

	citi := NewCiti()
	if !citi.Applicable(lines) {
		t.Log("File", file, "not recognized as citi")
		t.Fail()
	}
}
