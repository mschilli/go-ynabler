package ynabler

import (
	"go.uber.org/zap"
	"strings"
	"testing"
)

func TestChase(t *testing.T) {
	file := "test-data/chase.csv"

	lines, err := ReadFile(file)
	if err != nil {
		panic(err)
	}

	blog := zap.Must(zap.NewProduction())
	log := blog.Sugar()

	bofacc := NewChase()
	if !bofacc.Applicable(lines) {
		t.Log("File", file, "not recognized as chase")
		t.Fail()
	}

	out, err := bofacc.Process(log, lines)

	if err != nil {
		t.Log("Error processing chase file", err)
		t.Fail()
	}

	line := "02/27/2024,WHOLEFDS NOE 44378,,$33.59,"
	if !strings.Contains(out, line) {
		t.Log("line", line, "not found in", out)
		t.Fail()
	}
}
