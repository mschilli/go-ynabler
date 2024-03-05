package ynabler

import (
	"go.uber.org/zap"
	"strings"
	"testing"
)

func TestAmex(t *testing.T) {
	file := "test-data/amex.csv"

	lines, err := ReadFile(file)
	if err != nil {
		panic(err)
	}

	blog := zap.Must(zap.NewProduction())
	log := blog.Sugar()

	amex := NewAmex()
	if !amex.Applicable(lines) {
		t.Log("File", file, "not recognized as amex")
		t.Fail()
	}

	out, err := amex.Process(log, lines)

	if err != nil {
		t.Log("Error processing amex file", err)
		t.Fail()
	}

	if !strings.Contains(out, "02/29/2024,DIGITALOCEAN.COM NEW YORK CITY NY,,$6.00,") {
		t.Log("Digital Ocean line not found")
		t.Fail()
	}

	if !strings.Contains(out, "02/29/2024,MOBILE PAYMENT - THANK YOU,,,$603.38") {
		t.Log("Payment line not found")
		t.Fail()
	}
}
