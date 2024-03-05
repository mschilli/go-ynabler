package ynabler

import (
	"go.uber.org/zap"
	"strings"
	"testing"
)

func TestBofaCC(t *testing.T) {
	file := "test-data/bofacc.csv"

	lines, err := ReadFile(file)
	if err != nil {
		panic(err)
	}

	blog := zap.Must(zap.NewProduction())
	log := blog.Sugar()

	bofacc := NewBofaCC()
	if !bofacc.Applicable(lines) {
		t.Log("File", file, "not recognized as bofacc")
		t.Fail()
	}

	out, err := bofacc.Process(log, lines)

	if err != nil {
		t.Log("Error processing bofacc file", err)
		t.Fail()
	}

	if !strings.Contains(out, "03/02/2024,THE MISSION GROCERY SAN FANCISCO CA,,$13.36,") {
		t.Log("line not found")
		t.Fail()
	}
}
