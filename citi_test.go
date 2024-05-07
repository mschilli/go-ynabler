package ynabler

import (
	"go.uber.org/zap"
	"strings"
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

	blog := zap.Must(zap.NewProduction())
	log := blog.Sugar()
	out, err := citi.Process(log, lines)

	if err != nil {
		t.Log("File", file, "process error")
		t.Fail()
	}

	csv := strings.Split(out, "\n")

	exp := `02/27/2024,MANILA ORIENTAL MARKET DALY CITY CA,`
	if !strings.Contains(csv[3], exp) {
		t.Log("Line", csv[3], "Exp", exp)
		t.Fail()
	}

	exp = `02/28/2024,REFUND RHUBARB,,,$10.00`
	if !strings.Contains(csv[4], exp) {
		t.Log("Line", csv[4], "Exp", exp)
		t.Fail()
	}
}
